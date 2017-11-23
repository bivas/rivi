# Condition

The entire `condition` section is optional - you can run all rules all the time and see if it helps :smile:

- `if-labeled` - apply the rule if the issue has any of the provided labels
- `skip-if-labeled` - skip rule processing if issue has any of the provided labels
- `files`
    * `patterns` - [pattern](https://golang.org/s/re2syntax) matching the pull request file list (any of the patterns)
    * `extensions` - which file extension to match on pull request file list (must start with a dot [`.`])
- `ref`
    * `patterns` - [pattern](https://golang.org/s/re2syntax) matching the pull request ref name
    * `match` - matches the pull request ref name
- `title`
    * `starts-with` - issue title has a prefix
    * `ends-with` - issue title has a suffix
    * `patterns` - [pattern](https://golang.org/s/re2syntax) matching issue title (any of the patterns)
- `description`
    * `starts-with` - issue description has a prefix
    * `ends-with` - issue description has a suffix
    * `patterns` - [pattern](https://golang.org/s/re2syntax) matching issue description (any of the patterns)
- `comments`
    * `count` - number of comments for issue (supported operators: `==`, `>`, `<`, `>=`, `<=`)
- `order` - apply order hint to a rule. All rules are given order index **0**.
   **Important**: This will not place a rule in the exact position, but can assist in re-order rules.

## Example
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