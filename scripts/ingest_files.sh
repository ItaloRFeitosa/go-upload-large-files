#!/bin/bash

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <size_of_each_file_in_MB>"
    exit 1
fi

file_size=$1
max_total_size=$((10 * 1024)) # 10GB in MB
num_files=$((max_total_size / file_size))

echo "Number of files to upload: $num_files"

for i in $(seq 1 $num_files); do
    (
        timestamp=$(date +%s)
        dd if=/dev/zero bs=1M count=$file_size 2>/dev/null | \
        curl -X POST "http://localhost:8080/upload" \
             -H "Content-Type: multipart/form-data" \
             -F "file=@-;filename=large_file_${timestamp}_$i.bin"
        echo "Uploaded large_file_${timestamp}_$i.bin"
    ) &
    sleep 5
done

wait
echo "All uploads completed."