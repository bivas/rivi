# Label

## Goal

Ability to add or remove a label to/from an issue 

## Requirements

None

## Options

- `label` (optional) - the label to add
- `remove` (optional) - the label to remove

## Example
```yaml
rules:
    example:
      labeler:
        label: approved
        remove: rejected
```