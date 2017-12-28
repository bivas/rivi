# Running Rivi Locally

If you plan on running rivi on you local environment, please read though and make sure everything is set properly.

## Requirements

- Create a token with `repo` permissions
- Create a webhook and make sure the following are configured:

  - Select **content type** as `application/json`
  - Optionally, set a **secret** (this will be used by the bot to validate webhook content)
  - Register the following events
    - Pull request
    - Pull request review
    - Pull request review comment

  - If you have started rivi with several configuration files, you can set the hook URL to access each different file by passing `namespace` query param with the file name (without the `yaml` extension)
  Example: `http://rivi-url/?namespace=repo-x`

## Configuration File Structure

### Config Section

Configure the Git client for repository access and webhook validation
```yaml
config:
  provider: github
  token: my-very-secret-token
  secret: my-hook-secret-shhhhh
```

- `token` (required; unless set by env) - the client OAuth token the bot will connect with (and assign issues, add comments and lables)
- `provider` (optional) - which client to use for git connection - the bot tries to figure out which client to use automatically (currently only `github` is supported but others are on the way)
- `secret` (optional) - webhook secret to be used for content validation (recommended)

### Environment Variables

You can set the values for `token` and `secret` via environment variables:
`RIVI_CONFIG_TOKEN` and `RIVI_CONFIG_SECRET` respectively.

It is common to configure the bot by injecting environment variables via CI server.

## Executable
Rivi can be run as a bot which listens to incoming repository webhooks. This service must be internet facing to accept incoming requests (e.g. GitHub).
```
Usage: rivi bot [options] CONFIGURATION_FILE(S)...

	Starts rivi in bot mode and listen to incoming webhooks

Options:
	-port=8080				Listen on port
	-uri=/					URI path
```
### Example
```
$ rivi bot -port 9000 repo-x.yaml repo-y.yaml
```

## Docker

It is also possible to run Rivi as Docker container. Rivi's images are published to Docker Hub as `bivas/rivi`.

You should visit [bivas/rivi](https://hub.docker.com/r/bivas/rivi/) Docker Hub page and check for published [tags](https://hub.docker.com/r/bivas/rivi/tags/).

```
$ docker run --detach \
             --name rivi \
             --publish 8080:8080 \
             --env RIVI_CONFIG_TOKEN=<rivi oauth token> \
             --volume /path/to/config/files:/config \
             bivas/rivi rivi bot /config/repo-x.yaml
```