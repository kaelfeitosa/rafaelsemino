---
description: Rules for executing bash scripts on this Windows machine.
---
# Windows Execution Rules

When automating tasks or creating scripts on this host:

1. **Prefer `go run .` or Go binaries**: Golang is installed and works flawlessly cross-platform. Always prefer creating a Go CLI over trying to maintain Bash, Make, or Python scripts for complex or multi-OS tasks.
2. **Avoid `wsl` unless explicitly requested**: `wsl bash -c` can hang if the WSL subsystem is not initialized or expects interactive input.
3. **Use Git Bash for simple shell tasks**: If absolutely necessary to run a shell script through `run_command` in PowerShell, use:
   ```powershell
   & "C:\Program Files\Git\bin\bash.exe" scripts/script_name.sh
   ```
4. **Avoid Package Managers in Scripts**: `winget` and `scoop` often fail or require user interaction. Use Go dependencies or direct downloads.
