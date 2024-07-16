# Reverse Delimited Iterator

Reads records in reverse.

## Example

Given the file structure

```
record 1\n
record 2\n
record 3\n
```

The first 3 calls to `Iterator::next()` would produce the following 3 strings:

- record 3
- record 2
- record 1

Record delimiters are omitted.

Depending on the read size, multiple records might become available with a single read. Extra records are queued in memory and prioritized when available, avoiding file system reads.

That said, the worst case configuration is a read size of 1 byte.

## Traversal Algorithm

```
set last_index = file size - 1
set read_start = last_index - 1

set read_size = 1024
set working_buffer = ""
set read_buffer = [read_size]
read(read_start, len(read_buffer))


```
