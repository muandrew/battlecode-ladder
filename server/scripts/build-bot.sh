#!/bin/bash

DIR_BOTS=$1
BOT_UUID=$2

#todo more resilient

BOT_DIR=${DIR_BOTS}/${BOT_UUID}
BOT_WRK=${BOT_DIR}/workspace
cp -r bot-builder ${BOT_WRK}
unzip -u ${BOT_DIR}/source.jar -d ${BOT_WRK}/src
cd ${BOT_WRK}
./gradlew build
