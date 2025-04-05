# t5 filtering and modifiers

## filters

The filters should be implemented in the core module. We want certain structs that represent the filters. All filters implement an interface that returns a bool when presented items to filter. A filter function can then run these filters on a list of items.

For the filters there should be a defined flagset, so we can define those filters easily over the cli. A flagset also makes it easy to reuse them for different subcommands.

### Filter Flag Examples

```
--completed          # Show only completed tasks
--not-completed      # Show only non-completed tasks
--priority=A         # Filter by priority A
--priority=A,B,C     # Filter by priorities A, B, or C
--project=Home       # Filter tasks with +Home project
--context=@phone     # Filter tasks with @phone context
--due-today          # Filter tasks due today
--due-before=DATE    # Filter tasks due before DATE
--created-after=DATE # Filter tasks created after DATE
--has-tag=tag_name   # Filter tasks with specific tag
```

## modifiers

Modifiers are also implemented as structs. Those structs implement an interface that returns a modified item when presented an item to modify. A modifier function can then run these modifiers on a list of items.

Here we also want a flagset for the cli.

### Modifier Flag Examples

```
--complete           # Mark tasks as complete
--uncomplete         # Mark tasks as incomplete
--priority=A         # Set priority to A
--remove-priority    # Remove priority
--add-project=Home   # Add +Home project to tasks
--remove-project=Home # Remove +Home project from tasks
--add-context=@phone # Add @phone context to tasks
--remove-context=@phone # Remove @phone context from tasks
--due=DATE           # Set due date
--remove-due         # Remove due date
--append=TEXT        # Append text to task
--prepend=TEXT       # Prepend text to task
```