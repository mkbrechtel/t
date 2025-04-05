# Command Line Interface (CLI)

## Overview

This document describes the architecture of the CLI layer in the t5 application, emphasizing a clear separation of concerns between the command-line interface and the core application logic.

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