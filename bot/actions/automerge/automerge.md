# Automatic Merge

## Goal

Ability to automatically merge a pull request

Assignees must use one of the approval phrases (_Approve_, _Approved_, _LGTM_, _Looks good to me_ - case insensitive).

**Note** If others comment with one of the approval phrases, it will not count as approval

## Requirements

None

## Options

- `strategy` (optional) - which strategy to use when merging. can be `merge`, `squash` or `rebase` (default: `merge`)
- `require` (optional) - the number of approvals required (default: **0** meaning require all assignees to approve)
- `label` (optional) - label the issue instead of merging (dry-run)

## Example
```yaml
rules:
    example:
      automerge:
        require: 2
        strategy: squash
```
or
```yaml
rules:
    example:
      automerge:
        label: approved
```