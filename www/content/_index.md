---
title: Builder
description: "Builder is a tool that implements the builder specification. It is the specification to build and deploy software projects with agents. It focuses on speed and simplicity."
cascade:
  - _target:
      kind: page
    type: docs
  - _target:
      kind: section
    type: docs
---

# Helm Charts

You can add the Builder Helm repository to your local Helm configuration and search for available charts.

```bash
helm repo add builder https://zeiss.github.io/builder
helm repo update
helm search repo builder
```
