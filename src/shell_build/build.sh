#!/bin/sh

time=`date '+%Y-%m-%d-%H:%M:%S'`
version=`git log -1 --oneline`
strip_version=${version// /-}
go build -ldflags "-X main.time='${time}' -X main.version='${strip_version}'"
