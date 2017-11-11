# Configuration File Structure

In your `.rivi.rules.yaml` file you can add the following `roles` section

## Roles Section

List of roles for selecting (login) users for assignment.
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
```