# Service Description

This document describes the service responsibilities and the key components affected by the PR.

## Overview
The service is responsible for providing a RESTful API for creating, reading, updating, and deleting events and expenses.

## Service Responsibilities
The service responsibilities include:
* Creating new events.
* Retrieving existing events.
* Updating existing events.
* Deleting existing events.
* Adding new expenses to events.
* Retrieving all expenses for an event.
* Updating existing expenses.
* Deleting existing expenses.
* Calculating balances for events.
* Calculating settlements based on event balances (newly introduced).

## Key Modules / Components
The key modules and components include:
* `EventRepository`: An abstraction for interacting with event data storage.
* `InMemoryEventRepository`: An implementation of the `EventRepository` interface.
* `ExpenseRepository`: An abstraction for interacting with expense data storage.
* `InMemoryExpenseRepository`: An implementation of the `ExpenseRepository` interface.
* `BalanceCalculator`: A component responsible for calculating balances for events.
* `SettlementCalculator`: A component responsible for calculating settlements based on event balances (newly introduced).
