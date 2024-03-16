#!/bin/bash -ue

go install github.com/otakakot/repository-dispatch@latest

repository-dispatch --app-id ${APP_ID} --client-payload ${CLIENT_PAYLOAD} --event-type ${EVENT_TYPE} --repository-name ${REPOSITORY_NAME} --repository-owner ${REPOSITORY_OWNER} --secret-key ${SECRET_KEY}
