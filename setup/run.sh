#!/bin/bash

source env.sh

# install deps
sudo apt install libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev libzstd-dev libgflags-dev

. download_binaries.sh
. setup_config.sh
