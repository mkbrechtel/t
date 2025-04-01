# t5 End-to-End Tests

This directory contains end-to-end tests for the t5 application. These tests verify that the application works correctly as a whole by executing the application's commands and checking the results.

## Test Design Constraints

1. **Pure Go Implementation**: Tests are written as Go unit tests that call the t5 binary directly. No shell scripts are used.

2. **Binary Execution**: Tests execute the compiled t5 binary (expected to be in the parent directory) rather than calling Go code directly.

3. **File-Based Verification**: Tests verify correct behavior by inspecting the content of files after command execution, not by inspecting internal state.

4. **Self-Contained**: Tests are self-contained and can be run with `go test` from this directory.

5. **Configuration**: Tests use the included t5config.yaml file for configuration to ensure consistent test environment.

6. **Clean Environment**: Tests should set up and clean up their environment, leaving no traces after execution.

## Running the Tests

1. Build the t5 binary in the parent directory:
   ```
   cd ..
   go build -o t5
   ```

2. Run the tests:
   ```
   cd test
   go test -v
   ```

## Test Coverage

The end-to-end tests cover the following functionality:

1. **Task Update**: Verifying that the `todo update` command correctly processes todo.txt files.
2. **Task Listing**: Verifying that the `list` command correctly displays tasks from the event store.
3. **Configuration**: Verifying that the `config` command correctly displays configuration settings.
4. **Event Replay**: Verifying that the event replay system correctly rebuilds state from events.
5. **Task Lifecycle**: Verifying the complete lifecycle of a task from creation to completion.
6. **Priority Handling**: Verifying that task priorities are correctly handled.

## Test Files

- `t5_test.go`: Contains all the end-to-end tests
- `t5config.yaml`: Configuration file used by the tests
- `todo.txt`: Sample todo list used as test input