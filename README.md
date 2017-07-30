[![Build Status](https://travis-ci.org/bivas/rivi.svg?branch=development)](https://travis-ci.org/bivas/rivi)
[![Go Report Card](https://goreportcard.com/badge/github.com/bivas/rivi)](https://goreportcard.com/report/github.com/bivas/rivi)
[![codecov](https://codecov.io/gh/bivas/rivi/branch/development/graph/badge.svg)](https://codecov.io/gh/bivas/rivi)

# rivi - Automate your review process

Managing a repository can require tedious administrative work, and as pull requests increase, it becomes even more complex. Today’s pull review flow lacks important visible information that leads to serious management issues. 

Rivi is an innovative bot that automates repository management. Forget about manually checking which module was modified, or which people are in charge of a pull review, Rivi will do it for you. 
Rivi enables automatic labeling with common parameters so that administrators can immediately understand their repository status with a quick glance. It also assigns relevant people to pull request reviews, allows to add comments, merges pull requests, sends triggers to relevant users and notifying them about issues that require prompt attention and more. 

With Rivi, developers can focus on the actual code base and less on administrative unambiguous actions made every day.  We are looking to add more automation features to make the repository management process seamless, and our highest priority is to ensure that Rivi lives up to the community standards by providing true value and efficiency.

## Usage
Rivi can be run as a service which listens to incoming repository webhooks. This service must be internet facing to accept incoming requests (e.g. GitHub).
```
Usage: rivi	server [options] CONFIGURATION_FILE(S)...

	Starts rivi in server mode and listen to incoming webhooks

Options:
	-port=8080				Listen on port (default: 8080)
	-uri=/					URI path (default: "/")
```
### Example
```
$ rivi server -port 9000 repo-x.yaml repo-y.yaml
```

### Docker

It is also possible to run Rivi as Docker container. Rivi's images are published to Docker Hub as `bivas/rivi`.

You should visit [bivas/rivi](https://hub.docker.com/r/bivas/rivi/) Docker Hub page and check for published [tags](https://hub.docker.com/r/bivas/rivi/tags/).

```
$ docker run --detach \
             --name rivi \
             --publish 8080:8080 \
             --env RIVI_CONFIG_TOKEN=<rivi oath token> \
             --volume /path/to/config/files:/config \
             bivas/rivi rivi -config /config/repo-x.yaml
```

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

# Configuration File Structure

## Config Section

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

## Roles Section

List of roles for selecting (login) users for assignment. 
```yaml
roles:
  admins:
      - user1
      - user2
  reviewers:
      - user3
      - user4
  testers:
      - user2
      - user4
```

## Rules Section

Configure rules to be processed by the bot on each issue event. Each rule may have several actions.
### Structure

```yaml
rules:
  rule-name1:
    <condition>
    <action-name>
    <action-name>
    <action-name>
    ...
  rule-name2:
    <condition>
    <action-name>
...
```

### Example
```yaml
rules:
  rule-name:
      condition:
        if-labeled:
          - label1
          - label2
        skip-if-labeled:
          - label3
        files:
          patterns: 
            - "docs/.*"
          extensions: 
            - ".go"
        title:
          starts-with: "BUGFIX"
          ends-with: "WIP"
          patterns:
            - ".* Bug( )?[0-9]{5} .*"
        description:
          starts-with: "test PR please ignore"
          ends-with: "don't review yet"
          patterns:
            - ".*depends on #[0-9]{1,5}.*"
        ref:
          match: "master"
          patterns:
            - "integration_v[0-9]{2}$"
        comments:
          count: ">10"
        order: 5
      commenter:
        comment: "We have a match!"
      labeler:
        label: ready-for-review
```
### Condition

The entire `condition` section is optional - you can run all rules all the time and see if it helps :smile:
- `if-labeled` - apply the rule if the issue has any of the provided labels
- `skip-if-labeled` - skip rule processing if issue has any of the provided labels
- `files`
  - `patterns` - [pattern](https://golang.org/s/re2syntax) matching the pull request file list (any of the patterns)
  - `extensions` - which file extension to match on pull request file list (must start with a dot [`.`])
- `ref`
  - `patterns` - [pattern](https://golang.org/s/re2syntax) matching the pull request ref name
  - `match` - matches the pull request ref name
- `title`
  - `starts-with` - issue title has a prefix
  - `ends-with` - issue title has a suffix
  - `patterns` - [pattern](https://golang.org/s/re2syntax) matching issue title (any of the patterns)
- `description`
  - `starts-with` - issue description has a prefix
  - `ends-with` - issue description has a suffix
  - `patterns` - [pattern](https://golang.org/s/re2syntax) matching issue description (any of the patterns)
- `comments`
  - `count` - number of comments for issue (supported operators: `==`, `>`, `<`, `>=`, `<=`)
- `order` - apply order hint to a rule. All rules are given order index **0**. 
**Important**: This will not place a rule in the exact position, but can assist in re-order rules. 

### Available Actions
- [`autoassign`](bot/actions/autoassign/autoassign.md) - Automatic assignment of issue reviewers
- [`automerge`](bot/actions/automerge/automerge.md) - Automatic merge for approved pull requests
- [`commenter`](bot/actions/commenter/commenter.md) - Add comment to an issue
- [`labeler`](bot/actions/labeler/labeler.md) - Add/Remove label to/from an issue
- [`sizing`](bot/actions/sizing/sizing.md) - Size a pull request
- [`trigger`](bot/actions/trigger/trigger.md) - Send HTTP triggers
- [`locker`](bot/actions/locker/locker.md) - Lock an issue

# Example Configuration

```yaml
config:
  provider: github
  token: my-very-secret-token
  secret: my-hook-secret-shhhhh

roles:
  admins:
      - user1
      - user2
  reviewers:
      - user3
      - user4
  testers:
      - user2
      - user4

rules:
  pr-size:
        sizing:
          xs:
            label: size/xs
            changed-files-threshold: 5
          s:
            label: size/s
            changed-files-threshold: 15
          default:
            label: pending-approval
            comment: "Your pull-request is too large for review"

  docs:
        condition:
          files:
            patterns: 
              - "docs/.*"
        labeler:
          label: documentation

  assignment:
        condition:
          skip-if-labeled:
            - pending-approval
        autoassign:
          roles:
            - admins
            - reviewers
```
