import os
import re

entities_dirs = [
    "acervo/entities/actions",
    "acervo/entities/works",
    "acervo/entities/agents"
]

def get_entity_ids():
    ids = set()
    for d in entities_dirs:
        if not os.path.exists(d): continue
        for f in os.listdir(d):
            if f.endswith('.md'):
                ids.add(f[:-3])
    return ids

def get_dangling_images():
    with open('audit_output.txt', 'r') as f:
        lines = f.readlines()
    dangling = []
    for line in lines:
        if '[DANGLING IMAGE]' in line:
            parts = line.split(': ')
            if len(parts) > 1:
                dangling.append(parts[1].strip())
    return dangling

entity_ids = get_entity_ids()
dangling = get_dangling_images()

print(f"Total dangling images: {len(dangling)}")

# Try to find a match for each dangling image
for img in dangling:
    matched = False
    # Check if any entity ID is a prefix of the image name
    for eid in sorted(entity_ids, key=len, reverse=True): # sort by length to match longest prefix first
        if img.startswith(eid):
            print(f"MATCH: {img} -> Entity: {eid}")
            matched = True
            break
    if not matched:
        print(f"NO MATCH: {img}")
