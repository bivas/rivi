[![Build Status](https://travis-ci.org/bivas/rivi.svg?branch=development)](https://travis-ci.org/bivas/rivi)

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

## Requirements

- Create a token with `repo` persmissions
- Create a webhook and make sure the following are configured:


  - **content type** is `application/json`
  - Set a **secret** (this will be used by the bot to validate webhook content
  - Register the following event
    - Pull request
    - Pull request review
    - Pull request review comment

# Rule Structure

## Condition
```yaml
rule-name:
    condition:
      if-labeled:
        - label1
        - label2
      skip-if-labeled:
        - label3
      filter:
        pattern: "docs/.*"
        extension: ".go"
    commenter:
      comment: "We have a match!"
```
## Available Actions
- [`autoassign`](bot/actions/autoassign/autoassign.md) - Automatic assignment of issue reviewers
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
            pattern: "docs/.*"
        labeler:
          label: documentation

  assignment:
        condition:
          skip-if-labeled:
            - pending-approval
        autoassign:
          from-roles:
            - admins
            - reviewers
```
