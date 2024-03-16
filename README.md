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
Use in combination with [actions/create-github-app-token](https://github.com/marketplace/actions/create-github-app-token) as follows.

```yaml
jobs:
  dispatch:
    runs-on: ubuntu-22.04
    steps:
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.GITHUB_APP_ID }}
          private-key: ${{ secrets.GITHUB_APP_PRIVATE_KEY }}
          owner: owner          # need for other repository dispatch
          repositories: name    # need for other repository dispatch
      - name: Repository Dispatch
        uses: otakakot/repository-dispatch@v1.0.0
        with:
          token: ${{ steps.generate-token.outputs.token }}
          repository: owner/name
          event-type: event-type
          client-payload: "{\"payload\": \"xxxxxxxxxx\"}"
```
