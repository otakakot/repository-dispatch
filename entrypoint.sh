#!/bin/bash -ue

go install github.com/otakakot/repository-dispatch@latest

repository-dispatch --app-id ${APP_ID} --event-type ${EVENT_TYPE} --repository-owner ${REPOSITORY_OWNER} --repository-name ${REPOSITORY_NAME} --secret-key ${SECRET_KEY} --client-payload \'${CLIENT_PAYLOAD}\'
