# repository-dispatch

Create a repository dispatch event using the GitHub Apps token.

## Go Install

```shell
go install github.com/otakakot/repository-dispatch@latest
```

## Usage

```shell
repository-dispatch --token ${TOKEN} --event-type ${EVENT_TYPE} --repository-owner ${REPOSITORY_OWNER} --repository-name ${REPOSITORY_NAME} --secret-key "${SECRET_KEY}" --client-payload "${CLIENT_PAYLOAD}"
```

## Use GitHub Actions

Also available on GitHub Actions.

```yaml
jobs:
  dispatch:
    runs-on: ubuntu-latest
    steps:
      - name: Repository Dispatch
        uses: otakakot/repository-dispatch@v1
        with:
          token: token
          repository: owner/name
          event-type: event-type
          client-payload: "{\"payload\": \"xxxxxxxxxx\"}"
```
