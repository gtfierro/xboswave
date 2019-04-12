#!/bin/bash

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

function grant_if_not_exists() {
    attester=$1
    subject=$2
    perms=$3
    namespace=$4
    resource=$5
    filename=$6

    echo "\n" | wv rtprove --subject $subject -o $filename "${perms}@${namespace}/${resource}"
    if [[ $? != 0 ]]; then
        echo "\n" | wv rtgrant --attester $attester --subject $subject -e 3y --indirections 0 "${perms}@${namespace}/${resource}"
        echo "\n" | wv rtprove --subject $subject -o $filename "${perms}@${namespace}/${resource}"
    else
        wv verify $filename
    fi
}


cd example

create_entity namespace namespace.ent
create_entity service service.ent
create_entity client client.ent

grant_if_not_exists namespace.ent service.ent GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA==:serve_grpc namespace.ent xbospb/Test/*  serviceproof.pem

rm att*.pem
