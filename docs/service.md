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
* Calculating balances for events.
* Handling expense-related logic.

## Key Modules / Components
The key modules and components include:
* `EventRepository`: An abstraction for interacting with event data storage.
* `ExpenseUseCase`: A component responsible for handling expense-related logic.
* `BalanceCalculator`: A component responsible for calculating balances for events.
