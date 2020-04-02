#!/usr/bin/bash

bash check_dependency.sh
if [ ! $? -eq 0 ]; then
    exit 1
fi

# if [ -z "$(ps -A | grep redis-server)" ]; then
if [ -z "$(tmux ls | grep redis)" ]; then
    echo "starting redis"
    tmux new -d -s redis 'cd ~/ && redis-server'
    echo "waiting a bit"
    sleep 10
    echo "continuing"
fi

if [ ! -z "$(tmux ls | grep bcl)" ]; then
    echo "killing existing bcl server"
    tmux kill-session -t bcl
fi
echo "starting bcl server"
tmux new -d -s bcl 'cd go/app && go run main.go'

echo "deploy complete"
