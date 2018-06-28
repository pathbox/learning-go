# export GOPATH := $(CURDIR)
# This is how we want to name the binary output
# OUTPUT=main

BINARY=cat

# These are the values we want to pass for Version and BuildTime
GITTAG=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.GitTag=${GITTAG} -X main.BuildTime=${BUILD_TIME}"

all:
	go build ${LDFLAGS} -o ${BINARY} main.go