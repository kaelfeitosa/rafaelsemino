import os
import yaml
import glob
from collections import defaultdict

entities_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities"

# Go through all .md files in entities
works = glob.glob(os.path.join(entities_dir, "**", "*.md"), recursive=True)

# We want to extract for each entity:
# - title
# - ID
# - Attachments (url, category, label)

data = {}

for filepath in works:
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    if content.startswith("---"):
        parts = content.split("---", 2)
        if len(parts) >= 3:
            frontmatter_str = parts[1]
            try:
                fm = yaml.safe_load(frontmatter_str)
                if fm and 'id' in fm:
                    entity_id = fm['id']
                    attachments = []
                    # Find all attachment keys
                    i = 1
                    while True:
                        url_key = f"attachment_{i}_url"
                        cat_key = f"attachment_{i}_category"
                        label_key = f"attachment_{i}_label"
                        if url_key in fm:
                            url = fm.get(url_key, '').replace('[[', '').replace(']]', '')
                            cat = fm.get(cat_key, 'unknown')
                            label = fm.get(label_key, 'None')
                            attachments.append({
                                "index": i,
                                "url": url,
                                "category": cat,
                                "label": label
                            })
                            i += 1
                        else:
                            break
                    data[entity_id] = {
                        "title": fm.get("title", ""),
                        "attachments": attachments
                    }
            except Exception as e:
                pass

with open(r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\image_report.txt", "w", encoding='utf-8') as out:
    for entity_id, info in data.items():
        if not info["attachments"]:
            continue
        out.write(f"=== {entity_id} | {info['title']} ===\n")
        for att in info["attachments"]:
            out.write(f"  [{att['index']}] {att['url']} ({att['category']}) - {att['label']}\n")
        out.write("\n")

print("Report generated at c:\\Users\\mkael\\progs\\html\\rafaelsemino\\acervo\\cli\\image_report.txt")
