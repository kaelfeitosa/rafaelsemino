import re
import random
from collections import defaultdict

report_file = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\image_report.txt"
html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"

available = defaultdict(list)
all_images = []

with open(report_file, 'r', encoding='utf-8') as f:
    current_entity = None
    for line in f:
        line = line.strip()
        if line.startswith("==="):
            current_entity = line.split("|")[0].strip("= ")
        elif line.startswith("["):
            try:
                parts = line.split("] ", 1)[1]
                path_part, rest = parts.split(" (", 1)
                category = rest.split(")")[0]
                available[current_entity].append({"path": path_part, "category": category})
                all_images.append({"entity": current_entity, "path": path_part, "category": category})
            except Exception:
                pass

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

pattern = r"images/optimized/([^'\"]+\.webp)"
matches = re.finditer(pattern, html_content)

used_in_html = []
for m in matches:
    used_in_html.append(m.group(1))

used_in_html = list(dict.fromkeys(used_in_html))

assigned_images = set()

def pick_image_for(entity, context_hints=""):
    # get all available
    imgs = [img for img in available.get(entity, []) if img['path'] not in assigned_images]
    if not imgs:
        # fallback to something similar or agent-rafael-semino
        fallback_entities = ["agent-rafael-semino", "work-exu-nao-vem-hoje", "work-rastros-de-exu", "work-vao"]
        for fb in fallback_entities:
            imgs = [img for img in available.get(fb, []) if img['path'] not in assigned_images]
            if imgs:
                break
    
    if not imgs:
        # total fallback
        imgs = [img for img in all_images if img['path'] not in assigned_images]
        
    if not imgs:
        return None # out of images!
        
    # Prefer "registro" over "imprensa"
    preferred = [img for img in imgs if img['category'] == 'registro']
    if preferred:
        chosen = preferred[0]
    else:
        chosen = imgs[0]
        
    assigned_images.add(chosen['path'])
    return chosen['path']

# We map each unique path in HTML to a new image path.
mapping = {}
for old_path in used_in_html:
    entity = old_path.split("/")[0]
    new_path = pick_image_for(entity, old_path)
    if new_path:
        # The new path from report is like `entity/image.jpeg`
        # We need to construct the webp path: `entity/image.webp`
        # But wait, the build-assets script takes `images/optimized/...webp` 
        # and looks for `media/images/...jpeg`.
        # Are there subdirectories in optimized? 
        # The builder flattens or keeps structure?
        # Actually in ASSETS.md it says: 
        # html: `<img src="images/optimized/work-exu-nao-vem-hoje-001.webp">`
        # Wait, the current html has `<img src="images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-001.webp">`
        # So we should use `images/optimized/work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-001.webp`.
        # Let's format the new path to match.
        new_webp = new_path.rsplit(".", 1)[0] + ".webp"
        mapping[old_path] = new_webp
    else:
        print(f"FAILED TO MAP: {old_path}")

print("Mapping:")
for k, v in mapping.items():
    print(f" {k} -> {v}")

# Apply replacement
# We replace exact string "images/optimized/OLD" with "images/optimized/NEW"
new_html = html_content
for old, new in mapping.items():
    old_full = f"images/optimized/{old}"
    new_full = f"images/optimized/{new}"
    new_html = new_html.replace(old_full, new_full)

with open(html_file, 'w', encoding='utf-8') as f:
    f.write(new_html)
    
print("Updated index.html directly!")
