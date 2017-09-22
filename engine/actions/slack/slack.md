# Slack 

## Goal

Ability to send a Slack message 

## Requirements

`slack` section must be configured with valid api key (`api-key`)

\- or -

have `RIVI_SLACK_API_KEY` set as environment variable

## Options

### Slack Section
- `api-key` (required) - 
- `translator` (optional) - 

### Slack Action
- `message-template` (required) -
- `channel` (optional) - private message otherwise
- `notify` (optional) - `assignees` or `role` name.
  - If both `channel` and `notify` are empty - no message will be outputted

## Template Fields

- `SlackUser`
- `Number`
- `Title`
- `Owner`
- `Repo`
- `Origin`
  
## Example
```yaml
slack:
  api-key: my-slack-api-key
  translator:
    user1: slackUser1
    user2: anotherUser2

rules:
    example:
      slack:
        channel: dev
        notify: assignees
        message-tempalte: "@{{ .SlackUser }} You are assigned to issue {{ .Number }}"
```