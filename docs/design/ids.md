# UUID System

The UUID system in t5 uses UUIDv7 timestamps to generate time-based, monotonically increasing unique identifiers for tasks. We implement both short-form and long-form identifiers with the following features:

## Short and Long Form IDs

- **Long Form**: Standard UUID strings (36 characters) like `0192da75-c158-7d7f-be3c-d5b647bf7fa8`
- **Short Form**: Compact, URL-safe base64 encoding starting with 't' (typically ~22-25 characters) like `tI4JMLyHOGsqsq86FlAqspsrZt`

## Implementation Details

1. **Base64 Custom Encoding**:
   - Uses a modified base64 alphabet that shifts characters to ensure IDs start with 't'
   - Replaces special characters with safe alternatives to avoid parsing issues

2. **Ordering Preservation**:
   - Maintains the temporal ordering of UUIDs even in encoded form
   - Tasks created earlier will have "smaller" identifiers lexicographically

3. **Round-Trip Conversion**:
   - Any UUID can be converted to short form and back without data loss
   - The timestamp component is preserved with millisecond precision

## Task Identifier Handling

When ensuring task properties:
- If a task contains a valid UUID in either an `id:` or `uuid:` tag, it's preserved
- If `preferShortIDs` is true and the ID is a valid UUID, it gets converted to an `id:` tag with the short form
- Otherwise, it receives a `uuid:` tag with the standard UUID string
- Non-UUID identifiers in `id:` tags are preserved, and a new `uuid:` tag is automatically added
- All tasks are guaranteed to have a UUID for internal tracking and synchronization

This system enables efficient storage in todo.txt files while maintaining compatibility with external systems that expect standard UUIDs, while also allowing integration with systems that use their own ID formats.
