package cleaner

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var ErrNoObjectsToDelete = errors.New("no objects to delete")

func Start(ctx context.Context, s3client *s3.Client, bucket string) func() error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit the goroutine
				return
			case <-ticker.C:
				buckets, err := s3client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
					Bucket: &bucket,
				})
				if err != nil {
					slog.Error("Failed to list objects in S3 bucket", "error", err)
					continue
				}
				var objectsToDelete []string

				for _, object := range buckets.Contents {
					if time.Since(*object.LastModified) > 5*time.Second {
						objectsToDelete = append(objectsToDelete, *object.Key)
					}
				}

				if len(objectsToDelete) == 0 {
					continue
				}

				for _, key := range objectsToDelete {
					_, err := s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
						Bucket: &bucket,
						Key:    aws.String(key),
					})
					if err != nil {
						slog.Error("Failed to delete object from S3 bucket", "key", key, "error", err)
					}
				}

				slog.Info("Deleted objects from S3 bucket")
			}
		}
	}()

	return func() error {
		defer cancel()
		return teardown(ctx, s3client, bucket)
	}
}

func teardown(ctx context.Context, s3client *s3.Client, bucket string) error {
	for {
		if err := removeAll(ctx, s3client, bucket); err != nil {
			if errors.Is(err, ErrNoObjectsToDelete) {
				return nil
			}
			return err
		}
	}
}
func removeAll(ctx context.Context, s3client *s3.Client, bucket string) error {
	output, err := s3client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		slog.Error("Failed to list objects in S3 bucket", "error", err)
		return err
	}
	var objectsToDelete []string

	for _, object := range output.Contents {
		objectsToDelete = append(objectsToDelete, *object.Key)
	}

	if len(objectsToDelete) == 0 {
		return ErrNoObjectsToDelete
	}

	for _, key := range objectsToDelete {
		_, err := s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    aws.String(key),
		})
		if err != nil {
			slog.Error("Failed to delete object from S3 bucket", "key", key, "error", err)
		}
	}

	slog.Info("Deleted objects from S3 bucket")

	return nil
}
