#!/usr/bin/sh

if [ -z "$GOPATH" ];
then
	echo "GOPATH not set."
	exit 1
fi
IFS=':' read -ra GO <<< "$GOPATH"

DIR=$GO/src/github.com/muandrew
mkdir -p $DIR
echo "Created directory at:\n$DIR"
ln -s $PWD $DIR/battlecode-ladder
echo "Created sym link at:\n$DIR/battlecode-ladder\nto:\n$PWD"
