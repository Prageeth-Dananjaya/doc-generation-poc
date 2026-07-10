---
name: Test Agentic Workflow
on:
  pull_request:
    types: [opened, synchronize]
    branches: [stg]
permissions:
  contents: read
  issues: read
  pull-requests: read
  copilot-requests: write
tools:
  github:
    toolsets: [default]
---

# PR Summary Agent

Analyse the pull request and provide the summary.
