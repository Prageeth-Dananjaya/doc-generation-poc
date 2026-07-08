# Architecture

This document describes the overall architecture and how the PR changes fit into the system.

## Overview
The architecture consists of the following components:
* GitHub Actions workflow
* `update_docs.py` script
* Groq API
* Notion API
* GitHub API

## Components
The components include:
* GitHub Actions workflow: Triggers the `update_docs.py` script on PR events
* `update_docs.py` script: Runs the documentation orchestrator using the Groq API
* Groq API: Generates technical specifications based on PR changes
* Notion API: Fetches markdown content
* GitHub API: Fetches PR base and head refs

## Data Flow
The data flow includes:
* PR changes -> GitHub Actions workflow -> `update_docs.py` script -> Groq API -> Technical specifications
* PR changes -> GitHub Actions workflow -> `update_docs.py` script -> Notion API -> Markdown content
* PR changes -> GitHub Actions workflow -> `update_docs.py` script -> GitHub API -> PR base and head refs

## Dependencies
The dependencies include:
* Groq API
* Notion API
* GitHub API
