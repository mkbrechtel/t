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