# Architecture

This document describes the overall architecture and how the PR changes fit into the system.

## Overview
The architecture consists of the following components:
* Frontend application: A web application that sends requests to the backend API.
* Backend API: A RESTful API that provides endpoints for creating, reading, updating, and deleting events and expenses.
* EventRepository: An abstraction for interacting with event data storage.

## Components
The components include:
* Frontend application: Built using React and sends requests to the backend API.
* Backend API: Built using Node.js and Express.js, and provides RESTful endpoints for events and expenses.
* EventRepository: Implemented using an in-memory data store, but can be replaced with a database in the future.

## Data Flow
The data flow includes:
* Frontend application -> Backend API: The frontend application sends requests to the backend API to create, read, update, and delete events and expenses.
* Backend API -> EventRepository: The backend API interacts with the EventRepository to store and retrieve event data.
