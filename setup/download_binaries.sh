#!/bin/bash

source lib.sh

# Download the wavemq, waved and wv binaries and move them into the proper locations
if [ ! -f "$BIN_LOCATION/wavemq" ] || [ "$UPGRADE" = true ]; then
    $echo "${INFO}Downloading wavemq${NC}"
    curl -LO https://github.com/gtfierro/wavemq/releases/download/v1.4.0-alpha/wavemq
    chmod +x wavemq
    sudo cp wavemq $BIN_LOCATION/
else
    $echo "${OK}Already have wavemq${NC}"
fi

if [ ! -f "$BIN_LOCATION/waved" ] || [ "$UPGRADE" = true ]; then
    $echo "${INFO}Downloading waved${NC}"
    curl -LO https://github.com/gtfierro/wave/releases/download/v0.5-alpha/waved
    chmod +x waved
    sudo cp waved $BIN_LOCATION/
else
    $echo "${OK}Already have waved${NC}"
fi

if [ ! -f "$BIN_LOCATION/wv" ] || [ "$UPGRADE" = true ]; then
    $echo "${INFO}Downloading wv${NC}"
    curl -LO https://github.com/gtfierro/wave/releases/download/v0.5-alpha/wv
    chmod +x wv
    sudo cp wv $BIN_LOCATION/
else
    $echo "${OK}Already have wv${NC}"
fi
