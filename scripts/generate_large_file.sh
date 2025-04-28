#!/bin/bash

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <number_of_files> <size_of_each_file_in_MB>"
    exit 1
fi

num_files=$1
file_size=$2

for i in $(seq 1 $num_files); do
    timestamp=$(date +"%Y%m%d%H%M%S")
    dd if=/dev/zero of=./tests/testdata/large_file_${timestamp}_$i.bin bs=1M count=$file_size
done