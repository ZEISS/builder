# 👷‍♂️ Builder 

[![Test & Build](https://github.com/zeiss/builder/actions/workflows/main.yml/badge.svg)](https://github.com/zeiss/builder/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeiss/builder)](https://goreportcard.com/report/github.com/zeiss/builder)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Builder is a tool that implements the builder specification. It is the specification to build and deploy software projects with agents. It focuses on speed and simplicity.

🤹‍♀️ This is inspired by [quick](https://shopify.engineering/quick) from Shopify (♥️).

## Usage

```bash
builder help
```

## Authentication

The builder authenticates to the server using OpenID Connect. The authentication flow is handled by the server. It supports [dex](https://github.com/dexidp/dex) as an OpenID Connect provider.

To authenticate, the builder uses a [dex](https://github.com/dexidp/dex) client ID and secret. These are configured in the builder tool.

```bash
builder auth login --dex-client-id <client-id> --dex-client-secret <client-secret> --dex-client-url
```

This logs in the builder using the configured [dex](https://github.com/dexidp/dex) client ID and secret.

## Server

The server implements the deployment features of the builder specification.

## License

[Apache 2.0](/LICENSE)
