# Formats

## todo.txt
The todo.txt format is the primary data format for t5, following the specification from [todo.txt](http://todotxt.org/). This format uses a simple text-based structure that is human-readable and easily parseable.

Features:
- Task completion markers: `x` at the beginning of a line marks a completed task
- Priority markers: `(A)`, `(B)`, `(C)` at the beginning of a line indicate task priority
- Creation dates: `YYYY-MM-DD` format dates track when tasks were created
- Completion dates: Date when task was marked as complete (appears after the `x`)
- Projects: `+project` tags associate tasks with projects
- Contexts: `@context` tags specify where a task can be performed
- Key-value tags: `key:value` format for additional metadata

t5 extends the todo.txt format with additional features:
- UUID tracking: `uuid:` tags to uniquely identify tasks across systems
- Short IDs: `id:` tags provide shorter, human-readable identifiers
- Time tracking data: Metadata about time spent on tasks

## Markdown
t5 also supports interaction with Markdown files containing task lists. This enables seamless integration with note-taking apps and documentation systems.

The Markdown support:
- Converts Markdown task lists (`- [ ]` and `- [x]`) to todo.txt format internally
- Preserves todo.txt metadata in the Markdown format
- Supports bidirectional synchronization between Markdown files and the todo.txt database
- Allows embedding todo.txt-compatible metadata in Markdown task items
- Can extract tasks from Markdown headings for context awareness
- Supports both CommonMark and GitHub Flavored Markdown task lists

Example Markdown syntax:
```markdown
# Project Tasks

- [ ] Implement feature X +project @context due:2023-05-01
- [x] Fix bug #123 +project @bugfix id:tI4JiF07Mzhf8bKGG75Pq9
```
