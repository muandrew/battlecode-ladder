#!/bin/bash

DIR_BOTS=$1
DIR_MATCHES=$2
DIR_WORKERS=$3

WORKER_ID=$4
MATCH_UUID=$5

BOT_1_UUID=$6
BOT_1_NAME=$7

BOT_2_UUID=$8
BOT_2_NAME=$9

#optional
DIR_MAPS=${10}
MAP_NAME=${11}

#todo more resilient

BOT_1_WRK=${DIR_BOTS}/${BOT_1_UUID}/workspace/build/classes
BOT_2_WRK=${DIR_BOTS}/${BOT_2_UUID}/workspace/build/classes

MATCH_URL=${DIR_MATCHES}/${MATCH_UUID}

WKR_WRK=${DIR_WORKERS}/${WORKER_ID}/match

cd ${WKR_WRK}

./gradlew run \
-PteamA=${BOT_1_NAME} -PteamB=${BOT_2_NAME} \
-PteamAUrl=${BOT_1_WRK} -PteamBUrl=${BOT_2_WRK} \
-PmatchUrl=${MATCH_URL} \
-PmapsUrl=${DIR_MAPS} \
-Pmaps=${MAP_NAME} \
-Pwebsocket=$(( 8700 + WORKER_ID ))
