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

echo='echo -e'
RED='\033[1;31m'
BLUE='\033[1;34m'
GREEN='\033[1;32m'
CYAN='\033[1;36m'
YELLOW='\033[1;33m'
INFO=$YELLOW
PROMPT=$BLUE
ERROR=$RED
OK=$GREEN
NC='\033[0m' # No Color

if ! command_exists wv; then
    echo "Install wv"
    exit 1
fi

if ! command_exists docker; then
    echo "Install docker"
    exit 1
fi

check_var XBOS_DEMO_NAMESPACE_ENTITY
check_var XBOS_DEMO_ADMIN_ENTITY
check_var XBOS_DEMO_INGESTER_ENTITY
check_var XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY
check_var XBOS_DEMO_DRIVER_SYSTEM_MONITOR

curdir=$(pwd)

$echo "${INFO}Getting latest containers${NC}"
docker pull xbos/waved:latest
docker pull xbos/wavemq:latest
docker pull xbos/ingester:latest
docker pull xbos/driver_system_monitor

function setup_waved() {
    # set up waved
    docker kill xboswave-demo-setup-waved
    docker rm xboswave-demo-setup-waved
    OPUT=$(docker run -d --name xboswave-demo-setup-waved \
                -v ${curdir}/etc/waved/:/etc/waved/ \
                -p 910:910 \
                --restart always \
                xbos/waved:latest 2>&1)
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
                xbos/wavemq:latest 2>&1)
    echo $OPUT
}

function setup_influxdb() {
    #rm -rf influxdata
    mkdir -p influxdata
    docker kill xboswave-demo-setup-influxdb
    docker rm xboswave-demo-setup-influxdb
    OPUT=$(docker run -d -p 8086:8086 \
           -v ${curdir}/influxdata:/var/lib/influxdb \
           --name xboswave-demo-setup-influxdb \
           influxdb 2>&1)
    echo $OPUT
}

function setup_ingester() {
    # TODO: configure the ingester so it takes environment variables
    # TODO: enable some static configuration of the ingester?
    echo "setup"
    docker kill xboswave-demo-setup-ingester
    docker rm xboswave-demo-setup-ingester
    cp $XBOS_DEMO_INGESTER_ENTITY etc/ingester/.
    ssh-keygen -f etc/ingester/sshkey -N '' -t ed25519 -b 1024
    OPUT=$(docker run -d --name xboswave-demo-setup-ingester \
        --link=xboswave-demo-setup-wavemq \
        --link=xboswave-demo-setup-influxdb \
        -e XBOS_INGESTER_SITE_ROUTER="xboswave-demo-setup-wavemq:4516"\
        -e XBOS_INGESTER_INFLUXDB_ADDR="http://xboswave-demo-setup-influxdb:8086"\
        -e XBOS_INGESTER_ENTITY_FILE=/etc/ingester/${XBOS_DEMO_INGESTER_ENTITY} \
        -e WAVE_DEFAULT_ENTITY=/etc/ingester/${XBOS_DEMO_INGESTER_ENTITY} \
        -v ${curdir}/etc/ingester:/etc/ingester/ \
        -p 2222:2222 \
        --restart always \
        xbos/ingester:latest 2>&1)
    echo $OPUT
}

$echo "${INFO}Setting up WAVED, InfluxDB${NC}"
setup_waved
setup_influxdb

$echo "${INFO}Creating entities${NC}"
create_entity "XBOS_DEMO_NAMESPACE_ENTITY" $XBOS_DEMO_NAMESPACE_ENTITY
create_entity "XBOS_DEMO_ADMIN_ENTITY" $XBOS_DEMO_ADMIN_ENTITY
create_entity "XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY" $XBOS_DEMO_WAVEMQ_SITE_ROUTER_ENTITY
create_entity "XBOS_DEMO_DRIVER_SYSTEM_MONITOR" $XBOS_DEMO_DRIVER_SYSTEM_MONITOR
create_entity "XBOS_DEMO_INGESTER_ENTITY" $XBOS_DEMO_INGESTER_ENTITY

oput=$(wv inspect ${XBOS_DEMO_NAMESPACE_ENTITY})
echo $oput
namespace_hash=$(wv inspect ${XBOS_DEMO_NAMESPACE_ENTITY}  | grep Hash | awk '{print $2}')
check_var namespace_hash

$echo "${INFO}Granting permissions${NC}"
# create dots
echo "\n" | wv rtprove --subject $XBOS_DEMO_ADMIN_ENTITY -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
if [[ $? != 0 ]]; then
    echo "\n" | wv rtgrant --attester $XBOS_DEMO_NAMESPACE_ENTITY --subject $XBOS_DEMO_ADMIN_ENTITY -e 3y --indirections 5 "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
    echo "\n" | wv rtprove --subject $XBOS_DEMO_ADMIN_ENTITY -o test.pem "wavemq:publish,subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
else
    wv verify test.pem
fi

# ingester
echo "\n" | wv rtprove --subject $XBOS_DEMO_INGESTER_ENTITY -o ingester-proof.pem "wavemq:subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
if [[ $? != 0 ]]; then
    echo "\n" | wv rtgrant --attester $XBOS_DEMO_ADMIN_ENTITY --subject $XBOS_DEMO_INGESTER_ENTITY -e 3y --indirections 0 "wavemq:subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
    echo "\n" | wv rtprove --subject $XBOS_DEMO_INGESTER_ENTITY -o ingester-proof.pem "wavemq:subscribe@${XBOS_DEMO_NAMESPACE_ENTITY}/*"
else
    wv verify ingester-proof.pem
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

$echo "${INFO}Setting up WAVEMQ, Ingester${NC}"
setup_wavemq
setup_ingester


# setup drivers
# TODO: setup network

#$WAVEMQ_SITE_ROUTER
XBOS_NAMESPACE=$namespace_hash
BASE_RESOURCE=drivers/systemmonitor
WAVEMQ_SITE_ROUTER="xboswave-demo-setup-wavemq:4516"
WAVE_DEFAULT_ENTITY=$XBOS_DEMO_DRIVER_SYSTEM_MONITOR

$echo "${INFO}Starting drivers${NC}"

docker kill xbos-demo-driver-system-monitor
docker rm xbos-demo-driver-system-monitor
docker run -d --name xbos-demo-driver-system-monitor \
    --link=xboswave-demo-setup-wavemq \
    -e WAVEMQ_SITE_ROUTER=${WAVEMQ_SITE_ROUTER} \
    -e XBOS_NAMESPACE=${namespace_hash} \
    -e BASE_RESOURCE=drivers/systemmonitor \
    -e WAVE_DEFAULT_ENTITY=${XBOS_DEMO_DRIVER_SYSTEM_MONITOR} \
    --restart always \
    xbos/driver_system_monitor
docker cp $XBOS_DEMO_DRIVER_SYSTEM_MONITOR xbos-demo-driver-system-monitor:/app/.

$echo "${INFO}Adding archival requests${NC}"
#echo "ssh -p 2222 root@172.17.0.1"
cmd="add xbosproto/XBOS /plugins/system_monitor.so ${namespace_hash} *"
$echo "${OK}Password is 'demo'${NC}"
echo "$cmd" |  ssh -p 2222 -o "UserKnownHostsFile=/dev/null" root@172.17.0.1

$echo "${OK}Password is 'demo'${NC}"
echo "list" |  ssh -p 2222 -o "UserKnownHostsFile=/dev/null" root@172.17.0.1
