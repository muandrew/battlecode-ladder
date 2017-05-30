#!/bin/bash

WORKER_ID=$1

#todo more resilient

WKR_DIR=bl-data/worker/${WORKER_ID}
if [[ -d "$WKR_DIR" ]]; then
    rm -r ${WKR_DIR}
fi
WKR_WRK=${WKR_DIR}/match
mkdir -p ${WKR_DIR}
cp -r bot-builder ${WKR_WRK}
