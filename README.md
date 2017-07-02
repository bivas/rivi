[![Build Status](https://travis-ci.org/bivas/rivi.svg?branch=development)](https://travis-ci.org/bivas/rivi)
[![Go Report Card](https://goreportcard.com/badge/github.com/bivas/rivi)](https://goreportcard.com/report/github.com/bivas/rivi)

# rivi
Automate your review process with Rivi the review bot

## Usage
```
Usage of rivi:
  -config string
    	Bot configuration file
  -port int
    	Bot listening port (default 8080)
  -uri string
    	Bot URI path (default "/")
```
### Example
```
$ rivi -port 9000 -config repo-x.yaml
```

## Requirements

- Create a token with `repo` persmissions
- Create a webhook and make sure the following are configured:


  - Select **content type** as `application/json`
  - Optionally, set a **secret** (this will be used by the bot to validate webhook content)
  - Register the following events
    - Pull request
    - Pull request review
    - Pull request review comment

# Configuration File Structure

## Config Section

Configure the Git client for repository access and webhook validation
```yaml
config:
  provider: github
  token: my-very-secret-token
  secret: my-hook-secret-shhhhh 
```

- `token` (required) - the client OAuth token the bot will connect with (and assign issues, add comments and lables)
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
        filter:
          patterns: 
            - "docs/.*"
          extension: 
            - ".go"
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
- `filter`
  - `patterns` - [pattern](https://golang.org/s/re2syntax) matching the pull request file list (any of the patterns)
  - `extensions` - which file extension to match on pull request file list (must start with a dot [`.`])
- `order` - apply order hint to a rule. All rules are given order index **0**. 
**Important**: This will not place a rule in the exact position, but can assist in re-order rules. 

### Available Actions
- [`autoassign`](bot/actions/autoassign/autoassign.md) - Automatic assignment of issue reviewers
- [`automerge`](bot/actions/automerge/automerge.md) - Automatic merge for approved pull requests
- [`commenter`](bot/actions/commenter/commenter.md) - Add comment to an issue
- [`labeler`](bot/actions/labeler/labeler.md) - Add label to an issue
- [`sizing`](bot/actions/sizing/sizing.md) - Size a pull request
- [`trigger`](bot/actions/trigger/trigger.md) - Send HTTP triggers

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
          filter:
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
