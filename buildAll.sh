#!/bin/sh

GOOS=linux go build -o spaghetti-analyzer.linux
GOOS=darwin go build -o spaghetti-analyzer.macos
GOOS=windows go build -o spaghetti-analyzer.windows
