#!/bin/bash

# Check if xvfb-run is available
if ! command -v xvfb-run &> /dev/null; then
    echo "xvfb-run not found. Running tests without xvfb..."
    go test "$@"
    exit $?
fi

# Run tests with xvfb
xvfb-run -s '-screen 0 640x480x24' go test "$@" 