package handler

import (
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/italorfeitosa/go-upload-large-files/pkg/httpresponse"
)

func Upload(s3Client *s3.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r, s3Client, bucket)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request, s3Client *s3.Client, bucket string) {
	requestID := uuid.New().String()
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to retrieve file from request", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	key := requestID + "/" + header.Filename
	partSize := int64(10 * 1024 * 1024) // 10 MB
	if header.Size <= partSize {        // 10 MB limit
		slog.Debug("Uploading file directly to S3", "filename", header.Filename, "request_id", requestID)
		_, err = s3Client.PutObject(r.Context(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &key,
			Body:   file,
		})
		if err != nil {
			slog.Error("Failed to upload file to S3", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		slog.Debug("Uploading file in parts to S3", "filename", header.Filename, "request_id", requestID)
		// upload concurrently for large files
		uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
			u.PartSize = partSize
		})
		_, err := uploader.Upload(r.Context(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			slog.Error("Failed to upload file to S3", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("X-Request-ID", requestID)
	httpresponse.JSON(w, http.StatusOK, JSON{
		"status":     "success",
		"request_id": requestID,
		"filename":   header.Filename,
		"object_key": key,
	})

	slog.Info("File uploaded successfully", "filename", header.Filename)
}
