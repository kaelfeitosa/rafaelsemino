---
description: Reflect on agent execution, identifying improvements to workflows, skills, or rules, and applying them.
---

# Meta-Reflection & Improvement

## 1. Analysis
Review the recent execution history, focusing on:
- `task.md` progress and blockers.
- Tool usage patterns (failures, retries, inefficiencies).
- User feedback and corrections.
- Gaps in `workflows/`, `skills/`, or `rules/`.

Identify concrete improvements:
- **Workflows**: Missing steps, unclear instructions, missing automation (turbo).
- **Skills**: Missing information, outdated docs, better examples needed.
- **Rules**: Contradictions, missing constraints.
- **Tools**: Opportunities for new scripts or CLI tools.

## 2. Plan
Propose specific changes to the `.agent/` directory or project `tools/`.
Describe the "before" and "after" for each change.

## 3. Execution (Apply Changes)
Perform the necessary file edits or creations.
- Use `write_to_file` or `replace_file_content` to update Agent docs.
- Use `run_command` to test new tools if applicable.

## 4. Persist
Summarize the reflection and actions taken.