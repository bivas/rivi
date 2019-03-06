# Changed to rivi rules file

## Example `rivi.yaml`

```yaml
rules:
    rivi:
        condition:
          files:
            patterns:
              - ".rivi.yaml"
        labeler:
          label: rivi-rules
```

**Note** The label `rivi-rules` must exists in the repository settings  

## Result

When there pull-request tries to change `.rivi.yaml` file:
<p><img src="rivi.before.png"/></p>

Rivi will add label `rivi-rules`:
<p><img src="rivi.after.png"/></p>