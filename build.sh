#!/usr/bin/env bash

export GOARCH=amd64
os_list="linux darwin windows"
rm -f build/*

cd checks
go get
for os in $os_list; do
  echo "Building check_kafka_connect for $os"
  GOOS=$os go build -o ../build/check_kafka_connect.${os}.amd64
done
cd ..

cd cloudwatch
go get
for os in $os_list; do
  echo "Building healthy_task_count for $os"
  GOOS=$os go build -o ../build/healthy_task_count.${os}.amd64
done
cd ..

echo "Done"
