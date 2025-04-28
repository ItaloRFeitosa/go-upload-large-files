package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

// CreateS3Client initializes an S3 client with custom configuration for MinIO compatibility or default S3
func CreateS3Client() (*s3.Client, error) {
	useMinio := os.Getenv("USE_MINIO") == "true"

	if useMinio {
		return createS3ClientForMinio()
	}

	// Default S3 client configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3ClientForMinio() (*s3.Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:9000" // Default MinIO endpoint
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY is not set")
	}

	secretKey := os.Getenv("S3_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("S3_SECRET_KEY is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), func(lo *config.LoadOptions) error {
		lo.Region = "us-east-1" // MinIO does not require a specific region
		lo.BaseEndpoint = endpoint
		lo.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
				Source:          "MinIO",
			}, nil
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	slog.Warn("Using MinIO S3 client", "endpoint", endpoint)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // MinIO uses path-style URLs
	})

	// Check connection with S3 by listing buckets
	_, err = client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
