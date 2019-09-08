#!/bin/bash

source lib.sh


# apply templates
envsubst < etc/waved/waved.toml > waved.toml
envsubst < etc/wavemq/wavemq.toml > wavemq.toml

update_or_replace_config $WAVED_CFG_LOCATION waved.toml etc/waved
update_or_replace_config $WAVEMQ_CFG_LOCATION wavemq.toml etc/wavemq

$echo "${OK}WAVED and WAVEMQ config files in position${NC}"

# setup systemd files
# apply templates
envsubst < systemd/wavemq.service > wavemq.service
envsubst < systemd/waved.service > waved.service

update_or_replace_config /etc/systemd/system waved.service systemd
update_or_replace_config /etc/systemd/system wavemq.service systemd

$echo "${OK}WAVED and WAVEMQ systemd files in position${NC}"

# cleanup
rm -f wavemq.service waved.service waved.toml wavemq.toml

# enable + start wavemq/waved

$echo "${OK}Starting waved and setting to start at boot${NC}"
sudo systemctl enable waved
sudo systemctl start waved

$echo "${OK}Starting wavemq and setting to start at boot${NC}"
sudo systemctl enable wavemq
sudo systemctl start wavemq
