#!/usr/bin/env sh

set -e -u

WORK_DIR="$PWD"

mkdir -p web
mkdir -p worker

docker run --rm -v "$WORK_DIR/web":/keys concourse/concourse \
  generate-key -t rsa -f /keys/session_signing_key

docker run --rm -v "$WORK_DIR/web":/keys concourse/concourse \
  generate-key -t ssh -f /keys/tsa_host_key

docker run --rm -v "$WORK_DIR/worker":/keys concourse/concourse \
  generate-key -t ssh -f /keys/worker_key

cp ./worker/worker_key.pub ./web/authorized_worker_keys
cp ./web/tsa_host_key.pub ./worker
