#!/bin/bash

set -e

version=v0.0.4
os=$(uname -s)
arch=amd64
url=https://github.com/eduhds/gspm/releases/download/$version/gspm-${os,}-$arch.tar.gz

if command -v curl &> /dev/null; then
    curl -L $url -o /tmp/gspm.tar.gz
elif command -v wget &> /dev/null; then
    wget -O /tmp/gspm.tar.gz $url
else
    echo "Please install curl or wget and try again."
    exit 1
fi

tar -C /tmp -xzf /tmp/gspm.tar.gz && \
    sudo mv /tmp/gspm /usr/local/bin && \
    sudo chmod +x /usr/local/bin/gspm

if [ $? -ne 0 ]; then
    echo "Failed to install gspm."
    echo "See https://github.com/eduhds/gspm for more information."
    exit 1
else
    echo 'gspm installed successfully!'
fi