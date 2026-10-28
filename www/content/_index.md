---
description: "Builder is a tool that implements the builder specification. It is the specification to build and deploy software projects with agents. It focuses on speed and simplicity."
cascade:
  - _target:
      kind: page
    type: docs
  - _target:
      kind: section
    type: docs
---

##  Getting Started

Builder is a tool that implements the builder specification. 

There are two components to builder.

- The `sparkles` CLI tool
- The `server` that runs the builder interface

🤹‍♀️ This is inspired by [quick](https://shopify.engineering/quick) from Shopify and [Magic](https://engineering.wealthsimple.com/from-prompt-to-url-the-magic-behind-magic) by WealthSimple.

## 🔨 Installation

Use Homebrew to install the `sparkles` CLI tool on macOS or Linux.

```
brew install zeiss/builder-tap/sparkles
```

> This installs the `sparkles` tool.

## 🎨 Customize 

You can customize Builder in the `.builder.yml` configuration file.

* 🖼️ [Sites](customization/sites)

## 📀 Helm Charts

You can add the Builder Helm repository to your local Helm configuration and search for available charts.

```bash
helm repo add builder https://zeiss.github.io/builder
helm repo update
helm search repo builder
```
