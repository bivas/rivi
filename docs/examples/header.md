# Check if copyright header is missing

## Example `rivi.yaml`

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
            not-patterns: 
              - "(?i)(copyright)"
      labeler:
        label: missing-copyright
```

**Note** The label `missing-copyright` must exists in the repository settings  