# Pull Request Sizing

## Goal

Ability to set label to indicate PR size (by files or total changes) 

There's a special rule named `default` which will act as fallback in case of no match found. This rule is optional and if none provided and none matched - no label will be apply.

The rule will re-evaluate the PR on each event and update the label if required.

## Requirements

All size labels must exist

## Options

- `label` (required) - the label to apply
- `changed-files-threshold` (optional) - Maximum number of changed files
- `changes-threshold` (optional) - Maximum number of changes
- `comment` (optional) - a comment to post to PR

## Example
```yaml
rules:
    example:
      sizing:
        xs:
          label: size/xs
          changed-files-threshold: 5
          changes-threshold: 10
        s:
          label: size/s
          changed-files-threshold: 10
        default:
          comment: "Your PR is way too large for review!"
```