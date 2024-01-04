# Upload Large Files

## Description
Project to understand/document/apply how to upload large files with golang without affect overall container/server memory

## Checklist
- [ ] Service should have at max 128mb of ram
- [ ] Limit network traffic to easy observation of resource consumption
- [ ] Prometheus and Grafana for Observability
- [ ] Use minio as object storage
- [ ] Script to generate large file
- [ ] Script to upload file
- [ ] Script to clean up (remove large files and containers)
- [ ] Unit test proxy logic
- [ ] upload large files to minio using S3 API
- [ ] load test upload with k6
- [ ] Create prom metric of upload throughput