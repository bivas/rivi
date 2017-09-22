# Automatic Assignment 

## Goal

Ability to assign users to an issue.

The action will assign users from matched roles provided that are available "spots". 

## Requirements

`roles` section must be configured with valid user login

## Options

- `roles` (optional) - select assignees from these roles (default: roles section)
- `require` (optional) - how many assignees are required to assign (default: **1**)

## Example
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

rules:
    example:
      autoassign:
        roles:
          - reviewers
          - testers
        require: 2
```