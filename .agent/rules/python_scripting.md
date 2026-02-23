---
description: Constraints for writing and executing transient Python scripts in the Acervo project
---

# Transient Python Scripting Rules

Python is frequently used by autonomous agents for generic data extraction, web scraping (like `instaloader` or `PyMuPDF`), or bulk renaming routines. However, because the Acervo database relies on strictly-typed Markdown YAML frontmatter and a compiled Go indexing engine, all Python interactions must adhere to explicit safety protocols to prevent database corruption.

## 1. Absolute Priority of the Go CLI 

**Do NOT use crude Python string-replacement (e.g. `content.replace("media_url", "path")`) or Regex scripts to mutate Markdown YAML frontmatter.**
While Python strings are fast, they are brittle and blind to YAML structural rules. 
If an Agent determines that bulk modifications are required inside `entities/` (e.g., adding a tag to 50 files, or renaming a key), the Agent MUST script a loop that invokes the native compiler:
`go run . ingest update <entity_id> [key=value]`

The Go tool natively Unmarshals and Marshals the YAML, ensuring structural integrity is preserved. Use Python only to *read* or *parse* data and dynamically generate `acervo ingest` bash/Powershell commands.

## 2. Fail-Safe Temporary Directory Cleanup

When extracting media or scraping web sources, scripts often build transient folders (e.g. `temp_downloads/`). Agents must ensure these folders are wiped before execution terminates.
If a scrape loop encounters a Rate Limit (`403 Forbidden`) or raises an unhandled exception, simple `os.rmdir()` commands at the bottom of the script will be bypassed, abandoning trash on the user's disk.

**Every Python extraction script must wrap temporary storage logic in robust `try...finally` or `atexit` handlers:**

```python
import os, shutil, atexit

TEMP_DIR = "temp_extraction"

def cleanup():
    if os.path.exists(TEMP_DIR):
        shutil.rmtree(TEMP_DIR)

atexit.register(cleanup)

os.makedirs(TEMP_DIR, exist_ok=True)
# Extraction logic executes here...
# If it crashes, atexit guarantees cleanup.
```

## 4. Helper Script Isolation (Go & Python)

To prevent package namespace collisions (especially `main redeclared` errors in the Go CLI) and root directory clutter:

- **Isolated Packages**: Never place diagnostic or one-off Go scripts in the root `acervo/cli/` directory. Use a subdirectory like `acervo/debug/` or run them as standalone files outside the main module path if possible.
- **Immediate Cleanup**: Any helper script (`.py`, `.go`, `.sh`) created to solve a specific sub-task MUST be deleted as soon as its objective is achieved.
- **Silent Tooling**: Prefere updating the native Go CLI capabilities (e.g., adding a flag to `validate`) over creating a new permanent helper script.
