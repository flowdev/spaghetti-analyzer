#!/bin/sh

GOOS=linux go build -o spaghetti-analyzer.linux
GOOS=darwin go build -o spaghetti-analyzer.darwin
GOOS=windows go build -o spaghetti-analyzer.windows
