#!/bin/sh

go build -o go-simple-expression-eval
./go-simple-expression-eval "$1"
rm go-simple-expression-eval