#!/bin/bash

WORKER_ID=$1
MATCH_UUID=$2

USER_UUID_1=$3
BOT_UUID_1=$4
BOT_1_NAME=$5

USER_UUID_2=$6
BOT_UUID_2=$7
BOT_2_NAME=$8

#todo more resilient

BOT_1_WRK=$PWD/bl-data/user/${USER_UUID_1}/bot/${BOT_UUID_1}/workspace/build/classes
BOT_2_WRK=$PWD/bl-data/user/${USER_UUID_2}/bot/${BOT_UUID_2}/workspace/build/classes

MATCH_URL=$PWD/bl-data/match/${MATCH_UUID}

WKR_WRK=$PWD/bl-data/worker/${WORKER_ID}/match

cd ${WKR_WRK}
./gradlew run -PteamA=${BOT_1_NAME} -PteamB=${BOT_2_NAME} -PteamAUrl=${BOT_1_WRK} -PteamBUrl=${BOT_2_WRK} -PmatchUrl=${MATCH_URL}
