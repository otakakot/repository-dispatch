#!/bin/bash -ue

go install github.com/otakakot/repository-dispatch@latest

repos=(`echo ${REPOSITORY} | tr '/' ' '`)

repository-dispatch --token ${TOKEN} --event-type ${EVENT_TYPE} --repository-owner ${repos[0]} --repository-name ${repos[1]} --client-payload "${CLIENT_PAYLOAD}"
