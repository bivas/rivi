[![Build Status](https://travis-ci.org/bivas/rivi.svg?branch=development)](https://travis-ci.org/bivas/rivi)
[![Go Report Card](https://goreportcard.com/badge/github.com/bivas/rivi)](https://goreportcard.com/report/github.com/bivas/rivi)
[![codecov](https://codecov.io/gh/bivas/rivi/branch/development/graph/badge.svg)](https://codecov.io/gh/bivas/rivi)

# rivi - Simplify your review process

Managing a repository can require tedious administrative work, and as pull requests increase, it becomes even more complex. Today’s review flow lacks important visible information that leads to serious management issues. 

Rivi is an innovative tool that automates repository management. Forget about manually checking which module was modified, or which people are in charge of a pull review, Rivi will do it for you. 
Rivi enables automatic labeling with common parameters so that maitainers can immediately understand their repository status with a quick glance. It also assigns relevant people to pull request reviews, allows to add comments, merges pull requests, sends triggers to relevant systems and notifying them about issues that require prompt attention and more. 

With Rivi, developers can focus on the actual code base and less on administrative unambiguous actions made every day.  We are looking to add more automation features to make the repository management process seamless, and our highest priority is to ensure that Rivi lives up to the community standards by providing true value and efficiency.

## Usage
Rivi is available as a Github Application. Find out more at [rivi-cm.org](http://rivi-cm.org)

If you wish to host rivi on your local environment, please follow the [installation guide](docs/installation.md)

## Configuration File Structure

Place a `.rivi.rules.yaml` file at the repository root directory to be processed.

This configuration file might have multiple sections (depending on your scenario), but the only required one is `rules` section

## Rules Section

Configure rules to be processed on each issue event. Each rule may have several actions.
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

Find out about available `condition` options [here](docs/condition.md)

### Available Actions
- [`autoassign`](engine/actions/autoassign/autoassign.md) - Automatic assignment of issue reviewers
- [`automerge`](engine/actions/automerge/automerge.md) - Automatic merge for approved pull requests
- [`commenter`](engine/actions/commenter/commenter.md) - Add comment to an issue
- [`labeler`](engine/actions/labeler/labeler.md) - Add/Remove label to/from an issue
- [`sizing`](engine/actions/sizing/sizing.md) - Size a pull request
- [`trigger`](engine/actions/trigger/trigger.md) - Send HTTP triggers
- [`locker`](engine/actions/locker/locker.md) - Lock an issue
- [`slack`](engine/actions/slack/slack.md) - Send Slack messages

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
