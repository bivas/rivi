# Check if copyright header is missing

## Example `rivi.rules.yaml`

```yaml
rules:
    copyright-header:
      condition:
        match-kind: all
        files:
          extensions:
            - ".go"
            - ".java"
            - ".scala"
        patch:
          hunk:
            starts-at: 1
            pattern: "^((?!Copyright)[\s\S])*$"
      labeler:
        label: missing-copyright
```

**Note** The label `missing-copyright` must exists in the repository settings  