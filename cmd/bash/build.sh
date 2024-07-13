#!/bin/bash
cd "$(dirname "$0")"/../..
docker build -t ethereum-tracker-app .