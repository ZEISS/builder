---
title: "Sites"
weight: 20
---

Sites are a deployment target specifically interacting with the Builder API.
A site is self-service, citizen-development internal hosting service.

## Example

```yaml {filename=".builder.yaml"}
sites:
  path: .
  name: fizzy-buzzy
  ignore:
    - .builder
```

## Structure

A site is defined by the following fields:

- `path`: The path to the site's source code (default: current directory).
- `name`: The name of the site.
- `ignore`: A list of files to ignore when deploying the site.
