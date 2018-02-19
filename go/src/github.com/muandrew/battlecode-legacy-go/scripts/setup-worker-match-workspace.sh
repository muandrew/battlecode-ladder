#!/bin/bash

DIR_WORKERS=$1
WORKER_ID=$2

#todo more resilient

WKR_DIR=${DIR_WORKERS}/${WORKER_ID}
if [[ -d "$WKR_DIR" ]]; then
    rm -r ${WKR_DIR}
fi
WKR_WRK=${WKR_DIR}/match
mkdir -p ${WKR_DIR}
cp -r bot-builder ${WKR_WRK}
