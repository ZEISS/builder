# 👷‍♂️ Builder 

[![Test & Build](https://github.com/zeiss/builder/actions/workflows/main.yml/badge.svg)](https://github.com/zeiss/builder/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeiss/builder)](https://goreportcard.com/report/github.com/zeiss/builder)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Builder is a tool that implements the builder specification. It is the specification to build and deploy software projects with agents. It focuses on speed and simplicity.

🤹‍♀️ This is inspired by [quick](https://shopify.engineering/quick) from Shopify and [Magic](https://engineering.wealthsimple.com/from-prompt-to-url-the-magic-behind-magic) by WealthSimple.

## Usage

* 👨‍🎨 Design
* 🏗️ Build
* 🚢 Ship

```bash
builder help
```

## Installation

```bash
brew install zeiss/builder-tap/builder
```

## Features

* 🏛️ Static site hosting
* 🧭 Spec-driven building
* ...

## 🗺️ Discovery 

The configuration for the authentication provider and 
the API can be served on a `.well-known/builder-configuration` endpoint.

```
GET /.well-known/builder-configuration

{
  "oidc_issuer": "http://builder.internal:5556/dex",
  "api_url": "http://builder.internal:5556"
}
```

This is a way to simplify the configuration.

## Authentication

The builder authenticates to the server using OpenID Connect. The authentication flow is handled by the server. It supports [dex](https://github.com/dexidp/dex) as an OpenID Connect provider.

To authenticate, the builder uses a [dex](https://github.com/dexidp/dex) client ID and secret. These are configured in the builder tool.

```bash
builder auth login --url <builder-server>
```

This logs you in to the builder server.

```bash
builder --help
```

There are a couple of basic features.

```bash
Usage:
  builder [command]

Available Commands:
  auth        Authenticate the builder (default: dex)
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize a new config
  sites       Manages sites
  task        Manage tasks
```

## Server

The server implements the deployment features of the builder specification.

Using [Helm](https://helm.sh) to run the builder on Kubernetes.

```bash
helm repo add builder https://zeiss.github.io/builder
helm repo update
helm search repo builder
```

## License

[Apache 2.0](/LICENSE)
