#!/bin/sh
WORKDIR=${WORKDIR-.}

if echo "$1" | grep -q "WORKDIR="; then
  export $1
  shift
fi

cd $WORKDIR
exec "$@"
