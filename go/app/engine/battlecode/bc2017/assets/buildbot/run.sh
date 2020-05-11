## dir structure
# workspace/... # a gradle build project 
# run.sh # this file
# source.jar # from user upload

# workspace should be a setup gradle project
sunzip-cli source.zip -ms 15 -mm 10240 -md 102400 -d workspace/src
pushd workspace
chmod +x gradlew
./gradlew build
popd
cp -r workspace/build/classes result
