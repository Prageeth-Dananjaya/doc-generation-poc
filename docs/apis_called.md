# APIs Called

This document enumerates internal and external APIs invoked by code changes in the PR.

## Overview
The PR introduces changes to the `update_docs.py` script, which invokes the Groq API.

## Internal API Calls
There are no internal API calls introduced by the PR.

## External API Calls
The external API calls include:
* Groq API: The script uses the Groq API to generate technical specifications based on the PR changes.
* Notion API: The script uses the Notion API to fetch markdown content.
* GitHub API: The script uses the GitHub API to fetch PR base and head refs.

## Data Contract / Payloads
The data contract and payloads for the external API calls include:
* Groq API: The script sends a request to the Groq API with a prompt to generate technical specifications.
* Notion API: The script sends a request to the Notion API with a page ID to fetch markdown content.
* GitHub API: The script sends a request to the GitHub API with a PR base and head refs to fetch the PR changes.
