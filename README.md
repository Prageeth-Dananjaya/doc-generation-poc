# Expense Splitter Backend

This repository will host a Go backend application for collaborative events where participants share expenses fairly.

## Goal

The service will help users determine each person's share of an event expense, even when one or more participants already paid part or all of the cost. The system should calculate balances so that participants with higher contributions are reimbursed appropriately.

## Core use cases

- Create a collaborative event
- Add participants to the event
- Record expenses paid by one or more people
- Calculate each participant's remaining share
- Determine who should receive or pay money to settle balances

## Architecture approach

The implementation will follow hexagonal architecture principles:

- Domain layer for business rules and core entities
- Application layer for use cases
- Adapters for HTTP, persistence, and external services
- Dependency inversion so the core is independent from infrastructure

## Proposed initial scope

- Event and participant management
- Expense recording
- Balance calculation for non-payers and payers
- REST API for the main flows
