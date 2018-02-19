#!/bin/bash

if [ -z "$GOPATH" ];
then
	echo "GOPATH not set."
	echo "If you want to use the default go workspace for your other projects use \$HOME/go for your GOPATH."
	exit 1
fi

#https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

WORKSPACE=$DIR/go
export GOPATH=$GOPATH:$WORKSPACE
export PATH=$PATH:$WORKSPACE/bin

