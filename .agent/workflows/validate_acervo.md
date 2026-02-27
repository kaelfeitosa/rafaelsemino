---
description: How to validate, reindex, and audit the Acervo
---
# Workflow: Validate and Reindex Acervo

Run this workflow whenever you make changes or want to ensure the Acervo is in a healthy state.

## 1. Syntactic Validation
Check if all entities follow the mandatory fields and date formats.
```bash
// turbo
cd acervo/cli && go run main.go validate
```

## 2. Relational Audit & DB Generation
The `reindex` command generates the `db.sqlite` cache and runs a mandatory link check.
```bash
// turbo
cd acervo/cli && go run main.go reindex
```

## 3. Full Graph Verification
Perform a deep audit of the entire graph relations and asset consistency.
- **Relational Audit**: Checks links between Agents, Works, and Actions.
- **Entity Asset Audit**: Verifies that images referenced in Markdown exist in the acervo.
- **Frontend Asset Audit**: Verifies that `images/optimized/` references in `index.html` have corresponding masters in the acervo.
```bash
// turbo
cd acervo/cli && go run main.go verify
```

## 4. Fixing Blockers
If any command fails, use the CLI update engine to patch the faulty Markdown entity:
```bash
// turbo
cd acervo/cli && go run main.go ingest update <id> [key=value...]
```
