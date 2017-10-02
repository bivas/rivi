# Trigger

## Goal

Ability to trigger HTTP endpoints 

The following JSON will be posted to request:
```
{
  "time": "<current timestamp>",
  "message": "Triggered by Rivi",
  "item": {
    "repository": "<issue repository>",
    "state": "<issue state>",
    "id": <issue id>,
    "title": "<issue title>"
  }
}
```

The following headers will be set to request:

- `X-Rivi-Event=trigger` 
- `User-Agent=Rivi-Agent/1.0`

## Requirements

None

## Options

- `endpoint` (required) - the target endpoint
- `method` (optional) - HTTP method (`GET` or `POST`) (default: **POST**)
- `headers` (optional) - key-values to add to request. Must start with `X-` otherwise will not be included
- `body` (optional) - Template for posted data
- `content-type` (optional) - Set `Content-Type` header (default: **`application/json`**)

## Example
```yaml
rules:
    example:
      trigger:
        endpoint: "http://example.com/hook"
```