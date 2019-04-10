#!/bin/bash

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

function check_var {
    if [[ "${!1}" == "" ]]
    then
        echo "Please set \$$1"
        exit 1
    fi
}

if ! command_exists wv; then
    echo "Install wv"
    exit 1
fi

check_var XBOS_DEMO_NAMESPACE_ENTITY
check_var XBOS_DEMO_ADMIN_ENTITY
check_var XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY
check_var XBOS_DEMO_DRIVER_SYSTEM_MONITOR

curdir=$(pwd)

# TODO pull containers

function setup_waved() {
    # set up waved
    docker kill xboswave-demo-setup-waved
    docker rm xboswave-demo-setup-waved
    OPUT=$(docker run -d --name xboswave-demo-setup-waved \
                -v ${curdir}/etc/waved/:/etc/waved/ \
                -p 910:910 \
                --restart always \
                waved:latest 2>&1)
    echo $OPUT
    sleep 2
}

function create_entity() {
    entityname=$1
    filename=$2
    # create entities
    if [ ! -f ${filename} ]; then
        echo "Creating $entityname"
        wv mke -o $filename -e 10y --nopassphrase
        if [[ $? != 0 ]]
        then
            echo "Could not create $entityname ${filename}"
            exit 1
        fi
    fi
}

function setup_wavemq() {
    cp xbos-demo-namespace-routeproof.pem etc/wavemq/.
    cp $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY etc/wavemq/.


    docker kill xboswave-demo-setup-wavemq
    docker rm xboswave-demo-setup-wavemq

    OPUT=$(docker run -d --name xboswave-demo-setup-wavemq \
                -v ${curdir}/etc/wavemq/:/etc/wavemq/ \
                -p 9516:4516 \
                --restart always \
                wavemq:latest 2>&1)
    echo $OPUT
}

setup_waved
create_entity "XBOS_DEMO_NAMESPACE_ENTITY" $XBOS_DEMO_NAMESPACE_ENTITY
create_entity "XBOS_DEMO_ADMIN_ENTITY" $XBOS_DEMO_ADMIN_ENTITY
create_entity "XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY" $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY
create_entity "XBOS_DEMO_DRIVER_SYSTEM_MONITOR" $XBOS_DEMO_DRIVER_SYSTEM_MONITOR

namespace_hash=$(wv inspect namespace.ent  | grep Hash | awk '{print $2}')

set -x
# create dots
echo "\n" | wv rtprove --subject $XBOS_DEMO_ADMIN_ENTITY -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
if [[ $? != 0 ]]; then
    echo "\n" | wv rtgrant --attester $XBOS_DEMO_NAMESPACE_ENTITY --subject $XBOS_DEMO_ADMIN_ENTITY -e 3y --indirections 5 "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
    echo "\n" | wv rtprove --subject $XBOS_DEMO_ADMIN_ENTITY -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
else
    wv verify test.pem
fi

echo "\n" | wv rtprove --subject $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY -o xbos-demo-namespace-routeproof.pem "wavemq:route@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
if [[ $? != 0 ]]; then
    echo "\n" | wv rtgrant --attester $XBOS_DEMO_NAMESPACE_ENTITY --subject $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY -e 3y --indirections 0 "wavemq:route@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
    echo "\n" | wv rtprove --subject $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY -o xbos-demo-namespace-routeproof.pem "wavemq:route@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
else
    wv verify xbos-demo-namespace-routeproof.pem
fi

echo "\n" | wv rtprove --subject $XBOS_DEMO_DRIVER_SYSTEM_MONITOR -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/drivers/systemmonitor/*"
if [[ $? != 0 ]]; then
    echo "\n" | wv rtgrant --attester $XBOS_DEMO_ADMIN_ENTITY --subject $XBOS_DEMO_DRIVER_SYSTEM_MONITOR -e 3y --indirections 0 "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/drivers/systemmonitor/*"
    echo "\n" | wv rtprove --subject $XBOS_DEMO_DRIVER_SYSTEM_MONITOR -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/drivers/systemmonitor/*"
else
    wv verify test.pem
fi

setup_wavemq


# setup drivers
# TODO: setup network

#$WAVEMQ_SITE_ROUTER
XBOS_NAMESPACE=$namespace_hash
BASE_RESOURCE=drivers/systemmonitor
WAVEMQ_SITE_ROUTER="xboswave-demo-setup-wavemq:4516"
WAVE_DEFAULT_ENTITY=$XBOS_DEMO_DRIVER_SYSTEM_MONITOR

docker kill xbos-demo-driver-system-monitor
docker rm xbos-demo-driver-system-monitor
docker run -d --name xbos-demo-driver-system-monitor \
    --link=xboswave-demo-setup-wavemq \
    -e WAVEMQ_SITE_ROUTER=${WAVEMQ_SITE_ROUTER} \
    -e XBOS_NAMESPACE=${namespace_hash} \
    -e BASE_RESOURCE=drivers/systemmonitor \
    -e WAVE_DEFAULT_ENTITY=${XBOS_DEMO_DRIVER_SYSTEM_MONITOR} \
    --restart always \
    sys
docker cp $XBOS_DEMO_DRIVER_SYSTEM_MONITOR xbos-demo-driver-system-monitor:/app/.
