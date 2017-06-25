[![Build Status](https://travis-ci.org/bivas/rivi.svg?branch=development)](https://travis-ci.org/bivas/rivi)
[![Coverage Status](https://coveralls.io/repos/github/bivas/rivi/badge.svg?branch=development)](https://coveralls.io/github/bivas/rivi?branch=development)

# rivi
Automate your review process with Rivi the review bot

# Example

```yaml
config:
  provider: github
  token: my-very-secret-token
  secret: my-secret-shhhhh

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
