import os
import glob
import yaml

works_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities\works"
report_file = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\categories_report.txt"

carousel_works = [
    "work-irreversivel",
    "work-santo-bordel-de-tiatira",
    "work-chega-de-falar-de-botas",
    "work-de-louco-todo-mundo-tem-um-pouco",
    "work-de-sucupira-a-asa-branca",
    "work-a-serpente",
    "work-sociedade-o-circo"
]

lines = []
for w in carousel_works:
    fpath = os.path.join(works_dir, f"{w}.md")
    if os.path.exists(fpath):
        with open(fpath, 'r', encoding='utf-8') as file:
            content = file.read()
        if content.startswith('---'):
            yaml_content = content.split('---')[1]
            try:
                data = yaml.safe_load(yaml_content)
                lines.append(f"\n--- {w} ---")
                
                # Check all attachments for images
                for i in range(1, 10):
                    atype = data.get(f"attachment_{i}_type")
                    if atype == "image":
                        url = data.get(f"attachment_{i}_url", "")
                        cat = data.get(f"attachment_{i}_category", "UNKNOWN")
                        lines.append(f"[{i}]: {url} -> {cat}")
            except Exception as e:
                lines.append(f"Error parsing {w}: {e}")

with open(report_file, "w", encoding="utf-8") as f:
    f.write("\n".join(lines))
