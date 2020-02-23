#!/bin/bash

# executable path
DIR="$( cd "$( dirname "$0"  )" && pwd  )"
BAMBOO_DIR=$(dirname $DIR)
EXE_PATH=$BAMBOO_DIR/bamboo


if [ -f "$EXE_PATH" ]; then
 nohup $EXE_PATH >/dev/null 2>&1 &
 echo 'Start bamboo successfully! Default value http://127.0.0.1:6020'
else
 echo 'Cannot find $EXE_PATH.'
 exit 1
fi
