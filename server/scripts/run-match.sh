#!/bin/bash

WORKER_ID=$1
MATCH_UUID=$2

BOT_1_UUID=$3
BOT_1_NAME=$4

BOT_2_UUID=$5
BOT_2_NAME=$6

#todo more resilient

BOT_1_WRK=$PWD/bl-data/bot/${BOT_1_UUID}/workspace/build/classes
BOT_2_WRK=$PWD/bl-data/bot/${BOT_2_UUID}/workspace/build/classes

MATCH_URL=$PWD/bl-data/match/${MATCH_UUID}

WKR_WRK=$PWD/bl-data/worker/${WORKER_ID}/match

cd ${WKR_WRK}
./gradlew run -PteamA=${BOT_1_NAME} -PteamB=${BOT_2_NAME} -PteamAUrl=${BOT_1_WRK} -PteamBUrl=${BOT_2_WRK} -PmatchUrl=${MATCH_URL}
