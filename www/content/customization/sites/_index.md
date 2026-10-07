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

## Usage

Create the site using the `builder` CLI:

```bash
builder sites create --url <server>
```

Deploy the site:
```bash
builder sites deploy --url <server>
```

This will deploy the site to the specified server and make it available on the `<name>.<server>`.