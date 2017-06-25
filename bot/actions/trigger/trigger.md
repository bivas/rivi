# Trigger

## Goal

Ability to trigger HTTP endpoints 

The following JSON will be posted to request:
```
{
  "time": "<current timestamp>",
  "message": "Triggered by Rivi Bot",
  "item": {
    "repository": "<issue repository>",
    "state": "<issue state>",
    "id": <issue id>,
    "title": "<issue title>"
  }
}
```

The following headers will be set to request:

- `X-RiviBot-Event=trigger` 
- `User-Agent=RiviBot-Agent/1.0`

## Requirements

None

## Options

- `endpoint` (required) - the target endpoint
- `method` (optional) - HTTP method (`GET` or `POST`) (default: **POST**)
- `headers` (optional) - key-values to add to request. Must start with `X-` otherwise will not be included

## Example
```yaml
rules:
    example:
      trigger:
        endpoint: "http://example.com/hook"
```