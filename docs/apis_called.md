# APIs Called

This document enumerates internal and external APIs invoked by code changes in the PR.

## Overview
The PR introduces changes to the frontend application, which invokes internal APIs to create, read, update, and delete events and expenses.

## Internal API Calls
The internal API calls include:
* `POST /api/events`: Creates a new event.
* `GET /api/events`: Retrieves a list of all events.
* `GET /api/events/{id}`: Retrieves an event by ID.
* `PUT /api/events/{id}`: Updates an event.
* `DELETE /api/events/{id}`: Deletes an event.
* `GET /api/events/{id}/balances`: Calculates the balances for an event.

## External API Calls
There are no external API calls introduced by the PR.

## Data Contract / Payloads
The data contract and payloads for the internal API calls include:
* `POST /api/events`: The payload includes the event data, such as name, participants, and expenses.
* `GET /api/events`: The payload is empty.
* `GET /api/events/{id}`: The payload includes the event ID.
* `PUT /api/events/{id}`: The payload includes the updated event data.
* `DELETE /api/events/{id}`: The payload includes the event ID.
* `GET /api/events/{id}/balances`: The payload includes the event ID.
