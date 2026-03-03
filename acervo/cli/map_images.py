import re
import random
from collections import defaultdict

# 1. Read available images
report_file = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\image_report.txt"
available = defaultdict(list)
all_images = []

with open(report_file, 'r', encoding='utf-8') as f:
    current_entity = None
    for line in f:
        line = line.strip()
        if line.startswith("==="):
            current_entity = line.split("|")[0].strip("= ")
        elif line.startswith("["):
            # format: [1] path/to/image.jpeg (category) - label
            try:
                parts = line.split("] ", 1)[1]
                path_part, rest = parts.split(" (", 1)
                category = rest.split(")")[0]
                available[current_entity].append({"path": path_part, "category": category})
                all_images.append({"entity": current_entity, "path": path_part, "category": category})
            except Exception:
                pass

# 2. Extract current images from index.html
html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

# find all images/optimized/....webp
# Regex to match image paths in src or background url
pattern = r"images/optimized/([^'\"]+\.webp)"
matches = re.finditer(pattern, html_content)

used_in_html = []
for m in matches:
    used_in_html.append(m.group(1)) # e.g. "agent-rafael-semino/agent-rafael-semino-009.webp"

used_in_html = list(dict.fromkeys(used_in_html)) # preserve order, remove exact duplicates in usage

# We need to map each currently used path to a NEW path.
# Strategy:
# - Find the entity base name from the old path (e.g. agent-rafael-semino)
# - Pick an available "registro" image from that entity that hasn't been used yet.
# - If no "registro" is left, pick any.
# - If the entity has no more images, pick from a pool of generally applicable ones from agent-rafael-semino, etc..
# But wait, what if the context is specific to the work?
# Let's try to map the work to the same work's image if possible.

assigned_images = set()

def pick_image(entity, prefer_category="registro"):
    imgs = [img for img in available.get(entity, []) if img['path'] not in assigned_images]
    if not imgs:
        return None
    # Prefer category
    preferred = [img for img in imgs if img['category'] == prefer_category]
    if preferred:
        img = preferred[0] # Just pick the first available
        assigned_images.add(img['path'])
        return img['path']
    # Fallback to any
    img = imgs[0]
    assigned_images.add(img['path'])
    return img['path']

mapping = {}
unmapped = []

# Special contexts mapping according to old path
for old_path in used_in_html:
    entity = old_path.split("/")[0]
    new_path = pick_image(entity)
    if new_path:
        mapping[old_path] = new_path
    else:
        unmapped.append(old_path)

print(f"Total slots to fill: {len(used_in_html)}")
print("--- MAPPING ---")
for old, new in mapping.items():
    print(f"{old} -> {new}")

print("\n--- UNMAPPED (Need fallback) ---")
for u in unmapped:
    print(u)
