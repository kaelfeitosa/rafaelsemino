import os
import yaml

base_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities"

for root, dirs, files in os.walk(base_dir):
    for f in files:
        if f.endswith(".md"):
            path = os.path.join(root, f)
            try:
                with open(path, "r", encoding="utf-8") as file:
                    content = file.read()
                    if "---" in content:
                        parts = content.split("---")
                        if len(parts) >= 3:
                            yaml.safe_load(parts[1])
            except Exception as e:
                print(f"CORRUPT FILE FOUND: {path}")
                print(f"ERROR: {e}")

print("Sweep finished.")
