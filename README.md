# Upload Large Files - WIP

## Description

Project to understand/document/apply how to upload large files with golang without affecting overall container/server memory.

## Checklist

- [x] Service should have at max 256MB of RAM
- [ ] Limit network traffic to facilitate observation of resource consumption
- [x] LGTM for Observability
- [x] Use MinIO as object storage
- [x] Script to ingest large files
- [x] Worker to clean up (remove files from bucket)
- [ ] Unit test upload logic
- [ ] Integration test for upload logic
- [ ] Upload large files to MinIO using S3 API:
  - [x] Version using form-data
  - [ ] Version using multipart reader
  - [ ] Version using resumable uploads with tusd
- [ ] Set up load test for upload
- [x] Add containers dashboard to measure memory and CPU saturation
- [ ] Add API dashboard to measure rate, latency, and error rate

To test the health route of the service, use the following command:

```bash
curl -X GET http://localhost:8080/health
```

## Run Project Locally

Follow these steps to run the project locally:

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/italorfeitosa/go-upload-large-files.git
   cd go-upload-large-files
   ```

2. **Set Up Environment Variables**:
   Create a `.env` file in the project root (if not already present) and configure the required environment variables. Example:

   ```env
   MINIO_ROOT_USER=minioadmin
   MINIO_ROOT_PASSWORD=minioadmin
   MINIO_ACCESS_KEY=VPP0fkoCyBZx8YU0QTjH
   MINIO_SECRET_KEY=iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J
   BUCKET_NAME=uploads
   ```

3. **Start the Services**:
   Use Docker Compose to build and start the services:

   ```bash
   docker-compose up --build
   ```

   > **Note**: A worker runs in the background to clean up the bucket every 5 seconds. This prevents the disk from filling up with uploaded files. Ensure this behavior aligns with your testing or production needs.

4. **Verify the Setup**:

   - Access the MinIO console at `http://localhost:9001` using the credentials from the `.env` file.
   - Test the health route:
     ```bash
     curl -X GET http://localhost:8080/health
     ```

5. **Upload Files From Generator**:
   Use the provided script to upload files:

   ```bash
   ./scripts/ingest_files.sh <size_of_each_file_in_MB>
   ```

   > **Note**: Be careful with the size of the files, it can be harmful for host, still working in a sweet spot to not affect host

6. **Stop the Services**:
   To stop and clean up the services, run:
   ```bash
   docker-compose down -v
   ```
