# Upload Large Files

## Description
Project to understand/document/apply how to upload large files with golang without affect overall container/server memory

## Checklist
- [] Service should have at max 128mb of ram
- [] Limit network traffic to easy observation of resource consumption
- [] Prometheus and Grafana for Observability
- [] Use minio as object storage
- [] Script to generate large file
- [] Script to upload file
- [] Script to clean up (remove large files and containers)
- [] Unit test proxy logic
- [] E2E test for whole upload feature
- [] proxy upload to minio using S3 API