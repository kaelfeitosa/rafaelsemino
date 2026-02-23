---
name: Manage Acervo
description: Instructions and rules for managing the event-centric data structure of the Acervo.
---
# Manage Acervo Skill

You are responsible for managing the `acervo/` data structure. This is an event-centric database stored purely as Markdown files with YAML frontmatter.

## Core Directives

1. **Source of Truth**: The Markdown files in `acervo/entities/` are the single source of truth.
2. **Never Edit Binaries**: You must NEVER modify, rename, or delete files inside `acervo/registros/`.
3. **Immutability of IDs**: Once an entity has an `id` (e.g., `event-2023-teatro`), you must NEVER change it.
4. **Explicit Relations (Wikilinks)**: Never infer relations. Only create relations explicitly. ALL relational fields (like `related_to`, `created_by`, `work`, `agent`, `event`) MUST be formatted as Obsidian Wikilinks, e.g., `"[[agent-rafael-semino]]"`.
5. **No Interpretations**: Do not write artistic narratives, restructure the career, or deduce importance. Organize facts, don't interpret them.
6. **No Direct DB Edits**: NEVER write directly to `acervo/db.sqlite`. It is a derived index. Always edit the `.md` files and run the reindex script.
7. **Frontend Independence (Headless CMS)**: The `acervo/` backend is strictly headless. Do NOT copy backend assets (like images) into the `frontend/` directory. Markdown `media_url`s must explicitly point to the physical backend path (e.g., `/acervo/media/images/`).
8. **Terminal CMS Engine**: NEVER manually copy templates or write YAML dictionaries by hand to create new entities. You MUST use the native Go CLI: `go run . ingest create <type> <slug> [key=value...]` and `go run . ingest update <id> [key=value...]`.
9. **Transient Script Cleanup**: Always permanently delete transient Python extraction scripts from the project root immediately after the compilation/indexing phases are successful.
10. **Zero-Ambiguity Curation**: Never use fallback or generic agent mappings (e.g., defaulting to `agent-rafael-semino`) if a granular reading of the source text can explicitly map the record to a specific Work, Event, or Participation context.
11. **Mandatory UTF-8 Encoding**: Any automation script (PowerShell, Python, etc.) that interacts with or creates Acervo `.md` files MUST explicitly enforce `utf-8` encoding to prevent corrupted accents (e.g., `Ã§Ã£o`).
12. **Semantic Media Naming**: External media downloaded into `acervo/media/images/` must ALWAYS be named identically to its parent record ID (e.g., `record-ccbnb-001.jpeg`), never arbitrary labels like `instagram-image`.
13. **URL Provenance**: When scraping data from the web, if a physical media download fails (e.g. rate limits) or the post serves as an external citation, YOU MUST natively inject the source `url:` property directly into the entity's YAML frontmatter. Do not mock or generate fake placeholder image files to trick the system.

## Golang Execution Rules

14. **Compiler Cache Expiration**: The Go compiler caches binaries aggressively. If you ever modify any `.go` source files (e.g., `ingester.go`, `validator.go`), you MUST subsequently execute `go clean -cache` before running `go run . [command]` to ensure the executable reflects your architectural changes.
15. **Polymorphic AST Validation**: The Acervo contains heterogenous entities (e.g., Works have `title`, Agents have `name`, Participations have neither). Any edits to the `validator.go` AST logic must dynamically switch constraints based on the `type:` property. Never apply monolithic structural validations (like demanding `title`) globally across the entire graph.
16. **Native Reindex Authority**: The command `go run . reindex` is the absolute authoritative gateway to compiling the `db.sqlite` cache. Transient Python scripts cannot validate the Acervo graph. You must definitively pass the native Go `reindex` sequence to consider any data ingestion finalized.

## Entity Types and Locations

- **Agent**: `acervo/entities/agents/agent-<id>.md` (people or institutions)
- **Work**: `acervo/entities/works/work-<id>.md` (conceptual works)
- **Event**: `acervo/entities/events/event-<id>.md` (happenings in time/space)
- **Participation**: `acervo/entities/participations/participation-<id>.md` (agent participates in event, optionally with a work)
- **Record**: `acervo/entities/records/record-<id>.md` (documentary evidence pointing to `acervo/registros/`)

## Routine Tasks

When asked to update the Acervo, follow these steps:
1. **Locate/Scan**: Look for new files or new information provided.
2. **Determine Entity Type**: Agent, Work, Event, Participation, or Record.
3. **Create/Update**: Switch to the `acervo/cli/` directory and run either `go run . ingest create <type> <slug> [key=value...]` or `go run . ingest update <id> [key=value...]` to manipulate the Entity safely.
4. **Verify & Reindex**: Run `go run . reindex`. The Go compiler natively executes strict syntax validation and global graph relational audits prior to generating `db.sqlite`. It will forcefully block the index compilation if any markdown files contain broken links or sparse frontmatter.

## Acervo Auditing Guidelines

**Objective:** Guarantee that the acervo follows the correct conceptual model, has no misclassified entities, has no duplicate or implicit relations, has no orphaned entities, does not mix fact/action/evidence, and is safe for AI automation.

**Overall Process:** The audit happens in 4 layers, *always in this order*: Entities -> Relations -> Records -> Career Narrative. Never skip steps.

### 1. Entity Audit (Structural Sanity)
**1.1 Type Validation:** For each `.md` file, ask "Is the type of this entity unambiguous?"
- `Agent`: Someone or an institution?
- `Work`: Conceptual and timeless?
- `Event`: Has a date and location?
- `Participation`: Describes a concrete action?
- `Record`: Documentary evidence?
**Errors:** Event described as work, Work as event, Participation as event, Clipping as event. 
**Rule:** If it doesn't clearly answer "what it is", the type is wrong.

**1.2 Unique Identity (IDs):** For each entity: Is the ID unique? Is it stable? Does it appear only once?
**Red flags:** Two IDs for the same thing, IDs with conflicting dates, IDs based on fragile names ("event-new"). 
**Rule:** IDs are keys, not titles.

**1.3 Mandatory Fields:** Verify minimum fields. Do not try to infer missing fields; mark them for human correction.
- `Agent` -> `id`, `name`
- `Work` -> `id`, `title`, `created_by`
- `Event` -> `id`, `title`, `date_start`, `location`
- `Participation` -> `id`, `agent`, `event`, `role`
- `Record` -> `id`, `related_to`, `format`

### 2. Relation Audit (Semantic Sanity)
*This is the most important part.*

**2.1 Permitted Relations:** Check if all relations follow these rules:
- `Agent` -> `created_by` -> `Work`
- `Agent` -> `participates_in` -> `Participation`
- `Participation` -> `takes_place_in` -> `Event`
- `Participation` -> `uses_work` -> `Work`
- `Record` -> `documents` -> `Participation/Event`
**Forbidden:** Agent directly to Event, Work directly to Event, Record directly to Work (unless clear exception), Record to Agent without context. 
**Rule:** If these exist, it's a structural error.

**2.2 Relation Duplicity:** "Is this action registered more than once?"
**Examples of duplication:** Two identical participations (same agent/event/role), a participation + an event describing the same action, a record linked to two different participations without justification. 
**Rule:** One concrete action = one Participation. If there are two, unite or eliminate.

**2.3 Implicit Relations (Danger):** Verify textual descriptions that say something not modeled.
**Example:** Text says "directed play X in festival Y" but no corresponding Participation exists. This is a severe error for AI.
**Rule:** Nothing important should exist only in text.

### 3. Record Audit (Proofs and Evidences)
**3.1 Orphaned Records:** List records without `related_to` and records pointing to non-existent entities.
**3.2 Misplaced Records:** Check what the record actually proves. (e.g., clipping linked to Work instead of Participation, photo of event linked to Agent). It must point to what it proves.
**3.3 Record Redundancy:** Keep one Record. Reference it in multiple contexts only if it makes sense (rare). (e.g., same PDF registered twice, same article link with different IDs).

### 4. Career Narrative Audit (Global Coherence)
*(Human + AI step)*
**4.1 Timeline:** Build automatic timeline (Events by date, Participations by year, Works by period). Look for unexplained gaps, impossible overlaps, inconsistent/future dates.
**4.2 Role Progression:** Check role coherence over time, out-of-context abrupt transitions, duplicated roles in the same period. This is about historical sanity, not artistic judgment.
**4.3 Documented Impact:** For each relevant Participation, does at least one Record exist? Is there a clipping when there should be? If not, mark as a documentary gap (not an error).

### 5. Support Tools
- Use `audit.sh` script for: participations without records, orphaned records.
- Use SQLite queries for: duplicate participations, events without participation.
- Use Obsidian Graph View for: detecting isolated nodes, suspicious clusters.

### 6. Golden Rule of Auditing
If an entity or relation forces you to "explain out loud" what it is, it is wrong. A good model explains itself.

### 7. Recommended Rhythm
- Micro-audit (automatic by AI): Weekly
- Complete structural audit: Quarterly
- Deep conceptual audit: Yearly

### Final Summary
1. Validate types
2. Validate IDs
3. Validate permitted relations
4. Eliminate duplications
5. Resolve orphans
6. Check timeline
7. Check career coherence

**TOTAL FOCUS ON SIMPLIFICATION, CORRECTNESS, AND RELEVANCE.**
