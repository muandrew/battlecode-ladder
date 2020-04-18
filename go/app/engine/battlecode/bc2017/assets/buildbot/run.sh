## dir structure
# workspace/... # a gradle build project 
# run.sh # this file
# source.jar # from user upload

# for all results
mkdir -p result
# workspace should be a setup gradle project
sunzip-cli source.zip -ms 15 -mm 10240 -md 102400 -d workspace/src >> result/log.txt 2>&1
pushd workspace
chmod +x gradlew
./gradlew build >> ../result/log.txt 2>&1
popd
cp -r workspace/build/classes result
# zip up all files and put it back at root
pushd result
zip ../result.zip -r .
popd