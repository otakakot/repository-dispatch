#!/bin/bash -ue

key=`cat <<EOF
${SECRET_KEY}
EOF
`

payload=`cat <<EOF | jq
${CLIENT_PAYLOAD}
EOF
`

go install github.com/otakakot/repository-dispatch@latest

repository-dispatch --app-id ${APP_ID} --event-type ${EVENT_TYPE} --repository-owner ${REPOSITORY_OWNER} --repository-name ${REPOSITORY_NAME} --secret-key ${key} --client-payload ${payload}
