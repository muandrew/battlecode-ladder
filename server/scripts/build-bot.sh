#!/bin/bash

BOT_UUID=$1

#todo more resilient

BOT_DIR=$PWD/bl-data/bot/${BOT_UUID}
BOT_WRK=${BOT_DIR}/workspace
cp -r bot-builder ${BOT_WRK}
unzip -u ${BOT_DIR}/source.jar -d ${BOT_WRK}/src
cd ${BOT_WRK}
./gradlew build
