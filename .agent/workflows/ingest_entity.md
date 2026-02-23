---
description: How to ingest a new entity or record into the Acervo
---
# Workflow: Ingest Entity into Acervo

Follow these steps when you need to add new information or a new file to the Acervo.

## 1. Determine Entity Type
Identify if the new information represents an Agent, Work, Event, Participation, or Record.

## 2. Execute CLI Creation
Do NOT manually create Markdown files or copy templates. You MUST use the Acervo CLI to guarantee structural integrity. 
```bash
// turbo
cd acervo/cli && go run . ingest create <type> <slug> [key=value...]
```
*Example:* `go run . ingest create event festival-de-teatro name="Festival de Teatro" date="2023-01-01"`

## 3. Link Relations (Wikilinks)
If you need to update relationships (e.g., adding `related_to` or `work`), you MUST use Obsidian Wikilinks.
You can do this via the CLI:
```bash
// turbo
cd acervo/cli && go run . ingest update <id> related_to="[[agent-rafael-semino]]"
```

## 4. Reindex & Verify
The `reindex` command automatically runs `validate` (schema check) and `audit` (relational graph check) before touching the database.
```bash
// turbo
cd acervo/cli && go run . reindex
```
If the command blocks the index generation, read the error output and use `go run . ingest update <id>` to fix the broken links.
