# Go Tooling Improvement Plan

This document outlines a plan to enhance the `acervo/cli` tooling, consolidate scripts, and integrate deeper audit capabilities directly into the main Go application.

## 1. Project Structure & Module Configuration

### Issue
- The CLI is located in `acervo/cli`, but the module is named `acervo`. This structure works if run from `acervo/cli`, but imports like `acervo/internal/...` rely on the module root being there.
- Python scripts are mixed with Go source code.

### Proposal
- Ensure `go.mod` correctly defines the module.
- Move Python scripts to a `scripts/` subdirectory or deprecate them in favor of Go commands.
- Add a `Makefile` to `acervo/cli` (or root) to standardize commands like build, test, audit, and reindex.

## 2. Audit Capability Expansion (`acervo verify`)

### Issue
- The current `acervo/internal/auditor` only checks for broken links and participations without records.
- It misses the deep structural checks we performed manually (mandatory fields, forbidden relations, etc.).

### Proposal
- Port the logic from `audit_deep.py` into `acervo/internal/auditor/auditor.go`.
- Add checks for:
    - **Mandatory Fields**: Ensure Events have locations, Works have creators, etc.
    - **Forbidden Relations**: Block Agent->Event and Work->Event direct links.
    - **Suspicious Links**: Flag Records linking directly to Agents.
    - **Duplicate Participations**: Detect multiple participations for the same (Agent, Event, Role) tuple.

## 3. Tool Consolidation

### Issue
- `find_corrupt.py`, `fix_names.py`, `patch_names.py`, `run_mass_ingest.py` exist as standalone scripts.

### Proposal
- Integrate useful functionality into the `acervo` CLI.
    - `acervo tools find-corrupt`: Port logic from `find_corrupt.py`.
    - `acervo tools fix-names`: Port logic from `fix_names.py`.
- Deprecate the Python versions once ported.

## 4. Implementation Steps

1.  **Refactor Auditor**:
    - Update `internal/auditor/auditor.go` to parse frontmatter more robustly and implement the new checks.
    - Add specific error types or a report structure.
2.  **Add Makefile**:
    - Create `acervo/cli/Makefile` with targets: `build`, `reindex`, `verify`, `audit`.
3.  **Clean up**:
    - Move python scripts to `acervo/cli/legacy_scripts/` or similar.

## 5. Example Makefile

```makefile
.PHONY: build reindex verify audit

build:
	go build -o acervo-cli main.go

reindex:
	go run main.go reindex

verify:
	go run main.go verify

audit: verify
```
