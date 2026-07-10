# Logic Flows

This document describes the main business and data flow logic impacted by the PR.

## Overview
The changes introduced by the PR involve adding frontend draft persistence and a settlement plan UI.

## Entry Points
The entry points of the logic flow include:
* The frontend application, which sends requests to the backend API to create, read, update, and delete events and expenses.
* The `App.tsx` file, which handles the frontend application logic.

## Control Flow
The control flow of the logic flow involves the following steps:
1. The frontend application sends a request to the backend API to create a new event.
2. The frontend application autosaves the event draft to local storage.
3. The frontend application loads the saved draft from local storage.
4. The frontend application calculates the settlements based on the event balances.
5. The frontend application displays the suggested transfers.

## Important Branches
The important branches in the logic flow include:
* The `computeSettlements` function, which calculates the settlements based on the event balances.
* The `saveDraft` and `loadDraft` functions, which handle the frontend draft persistence.

## Changes Introduced
The changes introduced by the PR involve adding frontend draft persistence and a settlement plan UI. These changes impact the control flow of the logic flow, adding new steps and branches to handle the additional functionality.
