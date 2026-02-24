---
description: How to ingest a new entity (Action, Work, Agent) into the Acervo
---
# Workflow: Ingest Entity into Acervo

Follow these steps when you need to add new information to the Acervo.

## 1. Determine Entity Type
Identify if the new information represents an **Action**, **Work**, or **Agent**.

## 2. Execute CLI Creation
Do NOT manually create Markdown files or copy templates. You MUST use the Acervo CLI to guarantee structural integrity. 

```bash
// turbo
cd acervo/cli && go run main.go ingest create <type> <slug> [key=value...]
```
*Example:* `go run main.go ingest create action festival-teatro title="Atuação em Festival" performed_by="[[agent-rafael-semino]]" my_role="Ator"`

## 3. Link Relations (Wikilinks)
Relational fields (`performed_by`, `work_id`) MUST use Obsidian Wikilinks.
You can update them via the CLI:
```bash
// turbo
cd acervo/cli && go run main.go ingest update <id> [key=value...]
```

## 4. Reindex & Verify
```bash
// turbo
cd acervo/cli && go run main.go reindex && go run main.go verify
```
If the command blocks, read the error output and use `update` to fix the broken links.
