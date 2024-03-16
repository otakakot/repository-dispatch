# repository-dispatch

Create a repository dispatch event using the GitHub Apps token.

## Go Install

```shell
go install github.com/otakakot/repository-dispatch@latest
```

## Usage

```shell
repository-dispatch --app-id ${APP_ID} --event-type ${EVENT_TYPE} --repository-owner ${REPOSITORY_OWNER} --repository-name ${REPOSITORY_NAME} --secret-key "${SECRET_KEY}" --client-payload "${CLIENT_PAYLOAD}"
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
          app-id: ${{ secrets.GITHUB_APP_ID }}
          secret-key: ${{ secrets.GITHUB_APP_PRIVATE_KEY }}
          repository-owner: ${{ env.REPOSITORY_OWNER }}
          repository-name: ${{ env.REPOSITORY_NAME }}
          event-type: event-type
          client-payload: "{\"payload\": \"xxxxxxxxxx\"}"
```
