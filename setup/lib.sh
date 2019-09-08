#!/bin/bash

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

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

check_var() {
    if [[ "${!1}" == "" ]]
    then
        echo "Please set \$$1"
        echo "May need to source environment.sh"
        exit 1
    fi
}

update_or_replace_config() {
    # destination dir
    dest=$1
    filename=$2
    # extension
    template_dir=$3

    if [ ! -f $dest/$filename ] || [ "$OVERRIDE_CONFIG" = true ]; then
        sudo mkdir -p $dest
        sudo cp $filename $dest/$filename
    else
        differ=$(diff -q $template_dir/$filename $dest/$filename)
        if [ "$differ" = true ]; then
            $echo "${INFO} ${dest}/${filename} has a different configuration. Set \$OVERRIDE_CONFIG to true in env.sh to override${NC}"
            exit 1
        fi
    fi
}

function create_entity() {
    entityname=$1
    # create entities
    filename=${entityname}.ent
    if [ ! -f ${filename} ]; then
        echo "Creating $entityname"
        wv mke -o $filename -e 10y --nopassphrase
        if [[ $? != 0 ]]
        then
            echo "Could not create ${entityname} entity"
            exit 1
        fi
    fi
}

command_exists curl

check_var BIN_LOCATION
check_var UPGRADE
check_var WAVED_CFG_LOCATION
check_var WAVED_STORAGE_LOCATION
check_var WAVEMQ_CFG_LOCATION
check_var WAVEMQ_STORAGE_LOCATION
check_var OVERRIDE_CONFIG
