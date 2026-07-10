# APIs Called

This document enumerates internal and external APIs invoked by code changes in the PR.

## Overview
The PR introduces changes to the frontend application, which invokes internal APIs to calculate settlements and display suggested transfers.

## Internal API Calls
The internal API calls include:
* `computeSettlements`: Calculates the settlements based on the event balances.
* `saveDraft`: Saves the event draft to local storage.
* `loadDraft`: Loads the saved draft from local storage.

## External API Calls
There are no external API calls introduced by the PR.

## Data Contract / Payloads
The data contract and payloads for the internal API calls include:
* `computeSettlements`: The payload includes the event balances.
* `saveDraft`: The payload includes the event draft.
* `loadDraft`: The payload includes the saved draft.
