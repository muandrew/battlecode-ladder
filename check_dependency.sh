#!/usr/bin/bash

exit_result=0

function check_dependency {
    command=$1
    if [ -z "$(which ${command})" ]; then
        echo "missing ${command}"
        exit_result=1
    fi
}

check_dependency 'go'
check_dependency 'unzip'
check_dependency 'tmux'

if [ $exit_result -eq 0 ]; then
    echo "dependencies satisfied."
fi

exit $exit_result
