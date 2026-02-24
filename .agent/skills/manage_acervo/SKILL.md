# Manage Acervo Skill

You are responsible for managing the `acervo/` data structure. This is an editorial portfolio database stored purely as Markdown files with YAML frontmatter, centered on Actions, Works, and Agents.

## Core Directives

1. **Editorial Point of View**: The system is NOT neutral. It assumes authorial intent and a specific narrative recorte.
2. **Action-Centric Core**: Every entry must represent something Rafael DID.
3. **Source of Truth**: The Markdown files in `acervo/entities/` are the single source of truth.
4. **Immutability of IDs**: Once an entity has an `id` (e.g., `action-2023-teatro`), you must NEVER change it.
5. **Explicit Relations (Wikilinks)**: ALL relational fields (`performed_by`, `work_id`) MUST be formatted as Obsidian Wikilinks, e.g., `"[[agent-rafael-semino]]"`.
6. **No Direct DB Edits**: NEVER write directly to `acervo/db.sqlite`. It is a derived index. Always edit the `.md` files and run the reindex/verify scripts.
7. **Terminal CMS Engine**: You MUST use the native Go CLI for all entity manipulations: `go run main.go ingest create ...` or `go run main.go ingest update ...`.
8. **UTF-8 Enforcement**: All scripts and operations MUST preserve UTF-8 encoding (e.g., preservation of accents).

## Entity Types and Locations

- **Agent**: `acervo/entities/agents/agent-<id>.md` (Pessoas ou Coletivos. Rafael is a `person`).
- **Work**: `acervo/entities/works/work-<id>.md` (Obras artísticas. Work does not act; actions act upon works).
- **Action**: `acervo/entities/actions/action-<id>.md` (The facts. Includes embedded `context` and `attachments`).

## Absolute Rules for Actions

- Every Action represents something Rafael DID.
- No Action exists without `my_role`.
- No Action "of a collective" exists without Rafael's involvement.
- Action titles should be "portfolio phrases" (e.g., "Atuação em Vão", "Professor no Projeto Abarca").

## Audit & Verification Routine

When managing the Acervo:
1. **Locate**: Scan for new data sources.
2. **Ingest**: Use the Go CLI to create/update entities.
3. **Validate**: Run `go run main.go validate` (syntactic/schema check).
4. **Verify**: Run `go run main.go verify` (graph integrity/broken link audit).
5. **Reindex**: Run `go run main.go reindex` to sync the SQLite database.

## Relações Válidas
- `Action.performed_by` -> `Agent`
- `Action.work_id` -> `Work`
*Nenhuma outra relação é permitida.*
