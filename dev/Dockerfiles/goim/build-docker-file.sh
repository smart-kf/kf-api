#!/bin/bash

## build goim2.
set -ex
wd=$(pwd)
cd /tmp && git clone --depth=1 git@github.com:smart-kf/goim2.git
cd goim2
go mod tidy
make build
make discovery-local-image
make build-image
cd $wd
rm -rf /tmp/goim2