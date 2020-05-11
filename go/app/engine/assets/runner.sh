# This file is meant to setup for a run,
# log the situation, and  package the results.

mkdir -p result

if [[ -f source.sh ]]; then
    source source.sh
fi

bash run.sh > result/log.txt 2>&1
RUN_RESULT=$?

# zip up all files and put it back at root
pushd result
zip ../result.zip -r .
popd

exit $RUN_RESULT
