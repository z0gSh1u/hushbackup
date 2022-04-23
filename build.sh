#!/bin/bash

HUSHBACKUPVERSION="1.0"

# Target: Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/hushbackup-"$HUSHBACKUPVERSION"-linux-amd64

# Target: Linux i386
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o bin/hushbackup-"$HUSHBACKUPVERSION"-linux-amd64
