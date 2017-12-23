# Setting SHA Status

## Goal

Ability to set a status to pull request submitted commit

## Requirements

None

## Options

- `description` (required) - status description to set
- `state` (optional) - which state to set. can be `failure`, `pending` or `success` (default: `failure`)

## Example
```yaml
rules:
    example:
      status:
        description: pull request doesn't follow guidelines
        state: failure
```
or
```yaml
rules:
    example:
      status:
        description: pull request is okay
        state: success
```