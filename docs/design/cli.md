# Command Line Interface (CLI)

## Overview

This document describes the architecture of the CLI layer in the t5 application, emphasizing a clear separation of concerns between the command-line interface and the core application logic.

## Command Line Interface Structure

The t5 command line interface follows a consistent "subject-verb-object" sentence structure:

- **Subject**: `t5` - The name of our tool, which acts as the subject of each command.
- **Verb**: The action to perform, such as `add`, `modify`, `list`, or `update`.
- **Object**: What the action operates on, such as `todo` or a specific task ID.

Examples:
```bash
# Adding a new todo
t5 add todo "Buy groceries @shopping"

# Listing todos with a filter
t5 list todo --context=work

# Updating todo properties
t5 update todo
```

For operations on specific tasks, use the task ID directly as the object:

```bash
# Modify a specific task using its ID
t5 modify tI4JMLyHOGsqsq86FlAqspsrZt --priority=A
t5 modify tI4JMLyHOGsqsq86FlAqspsrZt --add-context=home

# Complete a specific task
t5 modify tI4JMLyHOGsqsq86FlAqspsrZt --complete
```

The task IDs (like `tI4JMLyHOGsqsq86FlAqspsrZt`) are short-form UUIDs that uniquely identify each task. This structure makes commands intuitive and easy to remember, while enabling both general operations on collections of tasks and specific operations on individual tasks.

## Architecture Principles

The main principles of our CLI design are:

1. **Strict Separation of Concerns**:
   - `cli` package: Responsible ONLY for parsing command-line arguments, building configuration, and delegating to the core application
   - `core` module: Contains ALL business logic including filtering, task modification, and state management
   - CLI should NOT contain any business logic - it should only prepare configurations and call application code

2. **CLI Responsibilities**:
   - Parse command-line arguments
   - Translate user flags into proper configurations
   - Instantiate appropriate core interfaces
   - Call the appropriate core functions with these configurations
   - Format and display results to the user

3. **Core Responsibilities**:
   - Implement all business logic
   - Process data based on configurations provided by CLI
   - Maintain application state
   - Implement domain-specific interfaces

## Benefits of This Architecture

1. **Testability**: Core logic can be tested independently of CLI
2. **Reusability**: Core components can be reused across different interfaces
3. **Maintainability**: Clear separation of concerns makes the codebase easier to understand and maintain
4. **Flexibility**: The core can be used with different interfaces (CLI, API, GUI) without changes
5. **Single Responsibility**: Each package has a clear and focused purpose

## Implementation

The application follows these separation guidelines:

- `cli/*.go`: Command-line interface definition, argument parsing, and output formatting
- `core/*.go`: All business logic and data processing
- `main.go`: Entry point that initializes the CLI

## Future Improvements

1. **Flag Reuse**: Extract flag definitions into a shared component for reuse across commands
2. **Command Discovery**: Implement dynamic command discovery for better extensibility
3. **Better Error Handling**: Improve error messages and handling for CLI users