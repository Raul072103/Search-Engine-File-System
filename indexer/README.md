# Overview

This is the component that handles the initialization of the database

## Component requirements

### 1. Events Queue

- ensures a predefined set operations that can be applied on the database
- is thread-safe

## TODO

- [X] Implement **Events Queue** so a set of database operations can be decided
- [X] Implement **Batch Processor** which takes operations from the Events Queue and applies them
- [X] Implement **File database initializer**
- [X] Implement **file crawler** (part of File database initializer)