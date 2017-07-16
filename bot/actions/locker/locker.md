# Locker

## Goal

Ability to lock/unlock an issue 

## Requirements

None

## Options

- `state` (required) - sets the issue state. Can be `lock`, `unlock` or `change`

## Example
```yaml
rules:
    example:
      locker:
        state: lock
```