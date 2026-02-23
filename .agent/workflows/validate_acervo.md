---
description: How to validate, reindex, and audit the Acervo
---
# Workflow: Validate and Reindex Acervo

Run this workflow whenever you make bulk changes or want to ensure the Acervo is in a healthy state.

## 1. Unified Verification & Reindexing
You do not need to run `validate` or `audit` manually. The `reindex` command automatically chains these processes in a strict succession:
1. Syntactic Schema Validation
2. Relational Graph Auditing
3. SQLite DB Generation

```bash
// turbo
cd acervo/cli && go clean -cache && go run . reindex
```

## 2. Fixing Blockers
If `reindex` halts due to a missing relation or broken schema, use the CLI update engine to patch the faulty Markdown entity:
```bash
// turbo
cd acervo/cli && go run . ingest update <id> [key=value...]
```
*Note: Always use Obsidian Wikilinks if updating relational fields like `related_to`!*
