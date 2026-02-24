---
description: Absolute priority of the Go CLI for data manipulation in the Acervo
---
# Go CLI Priority Rule

The Go-based CLI is the authoritative engine for all data operations in the Acervo.

## 1. No Manual YAML Manipulation
**Do NOT** manually edit the YAML frontmatter of Markdown files if a CLI command can achieve the same result. The CLI ensures type safety and structural integrity.

## 2. Command Priority
- **Creation**: Always use `go run main.go ingest create ...`.
- **Modification**: Always use `go run main.go ingest update ...`.
- **Verification**: Always use `go run main.go validate` and `verify`.
- **Indexing**: Always use `go run main.go reindex`.

## 3. Transient Script Safety
If you create transient scripts for bulk operations, they **MUST** invoke the CLI commands listed above rather than writing directly to `.md` files. This prevents structural corruption and maintains the "Action-centric" constraints.

## 4. Immediate Cleanup
Any one-off script created to assist in a task MUST be deleted immediately after execution.
