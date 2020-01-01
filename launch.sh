#!/usr/bin/bash

#todo check if tmux && redis is installed

# if [ -z "$(ps -A | grep redis-server)" ]; then
if [ -z "$(tmux ls | grep redis)" ]; then
    echo "starting redis"
    tmux new -d -s redis 'cd ~/ && redis-server'
    echo "waiting a bit"
    sleep 10
    echo "continuing"
fi

echo "starting server"
tmux new -d -s bcl 'cd go/src/github.com/muandrew/battlecode-legacy-go && go run main.go'
