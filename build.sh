#!/usr/bin/env bash

USER="dazwilkin" # Or your GitHub username
REPO="akri-http" # Or your preferred GHCR repo
TAGS="$(git rev-parse HEAD)"

for IMAGE in "client" "server" "devices"
do
  docker build \
  --tag=ghcr.io/${USER}/${REPO}-${IMAGE}:${TAGS} \
  --file=./deployment/Dockerfile.${IMAGE} \
  .
done
