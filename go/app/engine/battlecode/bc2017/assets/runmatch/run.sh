## dir structure
# workspace/... # a gradle build project 
# bot0.zip # the bot build result
# bot1.zip
# run.sh # this file
# source.sh # any params that needed to be passed

# Things that should be sourced
# WORKER_ID
# BOT_0_NAME
# BOT_1_NAME

map_name() {
    ls -1 -t map | head -1
}

sunzip-cli bot0.zip -ms 15 -mm 10240 -md 102400 -d bot0
sunzip-cli bot1.zip -ms 15 -mm 10240 -md 102400 -d bot1

DIR_MAPS=""
MAP_NAME=""
#optional
if [ -d "map" ]; then
    DIR_MAPS=$PWD/map
    map_file_path="$(map_name)"
    MAP_NAME=$(basename $map_file_path .map17)
fi

BOT_0_DIR=$PWD/bot0/classes
BOT_1_DIR=$PWD/bot1/classes

#.bc17 will be appended
MATCH_OUTPUT=$PWD/result/replay

pushd workspace
chmod +x gradlew
./gradlew run \
-PteamA=${BOT_0_NAME} \
-PteamAUrl=${BOT_0_DIR} \
-PteamB=${BOT_1_NAME} \
-PteamBUrl=${BOT_1_DIR} \
-PmatchUrl=${MATCH_OUTPUT} \
-PmapsUrl=${DIR_MAPS} \
-Pmaps=${MAP_NAME} \
-Pwebsocket=$(( 8700 + WORKER_ID ))
popd
