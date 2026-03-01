---
description: How to ingest a new entity (Work, Agent) or Occurrence into the Acervo
---
# Workflow: Ingest Entity into Acervo

Follow these steps when you need to add new information to the Acervo.

## 1. Determine Entity Type
Identify if the new information represents a **Work**, an **Agent**, or an **Occurrence** (which must belong to a Work).
The old `Action` entity no longer exists. All events now exist as Occurrences nested inside Works.

## 2. Execute CLI Creation
Do NOT manually create Markdown files or copy templates. You MUST use the Acervo CLI to guarantee structural integrity. 

```bash
// turbo
cd acervo/cli && go run main.go ingest create <type> <slug> [key=value...]
```
*Example:* `go run main.go ingest create work festival-teatro title="Atuação em Festival" medium="teatro" role="Ator"`

## 3. Link Relations (Wikilinks)
Relational fields (e.g. `collaborators`) MUST match Agent IDs where possible.
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
