#! /bin/bash
TARGET_IMAGE=$1
docker build -t trigger-resource --build-arg TARGET_IMAGE="$TARGET_IMAGE" .
