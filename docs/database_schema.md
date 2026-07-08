# Database Schema

This document captures the current database schema and any schema-related changes introduced by the PR.

## Overview
The changes introduced by the PR do not directly impact the database schema. However, the introduction of frontend draft persistence implies a potential future change to the database schema, potentially involving the creation of dedicated tables or collections for event drafts.

## Tables / Collections
Although no direct database schema changes are introduced by the PR, the following entities are likely to be represented in the database:
* Events
* Expenses
* Participants
* Payments
* Event Drafts (potential future addition)

## Relationships
The relationships between these entities can be described as follows:
* An event has multiple expenses (one-to-many).
* An expense belongs to one event (many-to-one).
* A participant can be part of multiple events (many-to-many).
* A payment is made by a participant for a specific expense (many-to-one).
* An event draft is associated with one event (one-to-one).

## Changes Introduced
There are no direct database schema changes introduced by the PR. However, the introduction of frontend draft persistence implies a potential future change to the database schema.
