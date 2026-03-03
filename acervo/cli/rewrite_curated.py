import re

html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

# The user wants:
# 1. Keep the original hero image.
# Looking at the original index.html (before my previous script), the hero image was:
# Wait, the hero image is in CSS:
# .hero--index { background: url('images/optimized/agent-rafael-semino/agent-rafael-semino-009.webp') center/cover no-repeat; }
# I need to ensure this is preserved.

# We need 33 OTHER unique images for the rest of the site.
# Let's map out exactly what images to use for what slot.
# I will curate a precise dictionary mapping line number / old path to new path.

mapping = {
    # HERO (Line 526)
    "agent-rafael-semino/agent-rafael-semino-009.webp": "agent-rafael-semino/agent-rafael-semino-009.webp", # Keep original
    
    # PERFIL (Lines 1338, 1340) - Needs 2 images of Rafael Semino
    "agent-rafael-semino/agent-rafael-semino-001.webp": "agent-rafael-semino/agent-rafael-semino-008.webp", # Great portrait
    "agent-rafael-semino/agent-rafael-semino-002.webp": "agent-rafael-semino/agent-rafael-semino-005.webp", # Another good portrait
    
    # EXU NÃO VEM HOJE - Autoria (Lines 1449, 1451, 1453) - Needs 3 images
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-002.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-008.webp", 
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-003.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-010.webp",
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-004.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-012.webp",
    
    # VÃO (Lines 1490, 1492) - Needs 2 images
    "work-vao/work-vao-001.webp": "work-vao/work-vao-003.webp",
    "work-vao/work-vao-002.webp": "work-vao/work-vao-004.webp",
    
    # CONTOS DE EXU (Lines 1525, 1527) - Needs 2 images
    "work-contos-de-exu/work-contos-de-exu-001.webp": "work-contos-de-exu/work-contos-de-exu-001.webp",
    "work-contos-de-exu/work-contos-de-exu-002.webp": "work-contos-de-exu/work-contos-de-exu-002.webp",
    
    # EXU NÃO VEM HOJE - Atuação (Lines 1567, 1569, 1571) - Needs 3 DIFFERENT images
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-011.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-001.webp",
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-009.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-007.webp",
    "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-010.webp": "work-exu-nao-vem-hoje/work-exu-nao-vem-hoje-009.webp",
    
    # MIRAIRA / REISADO (Lines 1625, 1628) - Needs 2 images
    "work-rafael-miraira-reisado/work-rafael-miraira-reisado-001.webp": "work-rafael-miraira-reisado/work-rafael-miraira-reisado-001.webp",
    "work-rafael-miraira-reisado/work-rafael-miraira-reisado-002.webp": "work-rafael-miraira-reisado/work-rafael-miraira-reisado-003.webp",
    
    # MESTRES DO MUNDO (Lines 1655, 1659, 1663, 1667) - Needs 4 images
    "work-mestres-do-mundo/work-mestres-do-mundo-001.webp": "work-mestres-do-mundo/work-mestres-do-mundo-001.webp",
    "work-mestres-do-mundo/work-mestres-do-mundo-002.webp": "work-mestres-do-mundo/work-mestres-do-mundo-002.webp",
    "work-mestres-do-mundo/work-mestres-do-mundo-003.webp": "work-mestres-do-mundo/work-mestres-do-mundo-003.webp", # Actually doesn't exist? Wait, let's check report: 001.jpeg, 001.png, 002.jpeg, 004.jpeg.
    # Let's fix MESTRES DO MUNDO
    # Re-mapping everything directly:
}

# The above was a manual start, but since we have the exact HTML paths, let's read the report and assign them carefully.

report_file = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\image_report.txt"
available = []

with open(report_file, 'r', encoding='utf-8') as f:
    for line in f:
        line = line.strip()
        if line.startswith("["):
            parts = line.split("] ", 1)[1]
            path_part, rest = parts.split(" (", 1)
            category = rest.split(")")[0]
            # convert jpeg/png to webp for HTML mapping
            webp_path = path_part.rsplit(".", 1)[0] + ".webp"
            # Don't use hero image in the pool
            if webp_path != "agent-rafael-semino/agent-rafael-semino-009.webp":
                available.append({"path": webp_path, "category": category})

# To not repeat, we track used
used = set(["agent-rafael-semino/agent-rafael-semino-009.webp"])

# Let's write a targeted function to apply replacements directly
import sys

def replace_nth(string, sub, wanted, n):
    where = [m.start() for m in re.finditer(re.escape(sub), string)]
    if len(where) <= n:
        return string
    idx = where[n]
    return string[:idx] + string[idx:].replace(sub, wanted, 1)

def get_best_image(prefix, prefer_registro=True):
    # Try find matching prefix
    candidates = [img for img in available if img['path'].startswith(prefix) and img['path'] not in used]
    if not candidates:
        return None
        
    if prefer_registro:
        regs = [c for c in candidates if c['category'] == 'registro']
        if regs:
            chosen = regs[0]['path']
            used.add(chosen)
            return chosen
            
    chosen = candidates[0]['path']
    used.add(chosen)
    return chosen

# We need to replace exactly the `<focus-image src="images/optimized/...">` in order.
# But `background: url(...)` is also there.
# Let's parse all image sources first.
pattern = r"images/optimized/([^'\"]+\.webp)"
matches = list(re.finditer(pattern, html_content))

new_html = html_content
offset = 0

for m in matches:
    old_path = m.group(1)
    
    # Hardcode hero keep
    if "center/cover" in html_content[m.start()-50:m.end()+20] or old_path == "agent-rafael-semino/agent-rafael-semino-009.webp":
        continue # skip
        
    entity_prefix = old_path.split("/")[0]
    
    new_path = get_best_image(entity_prefix, prefer_registro=True)
    
    if not new_path:
        # Fallback to ANY unused image, preferably from agent-rafael-semino
        ag_candidates = [img for img in available if img['path'].startswith('agent-rafael-semino') and img['path'] not in used]
        if ag_candidates:
             new_path = ag_candidates[0]['path']
        else:
             # Just any unused
             any_c = [img for img in available if img['path'] not in used]
             if any_c:
                 new_path = any_c[0]['path']
             else:
                 print(f"CRITICAL: NO IMAGES LEFT FOR {old_path}")
                 continue
                 
    used.add(new_path)
    
    # Replace exactly at this position
    start = m.start(1) + offset
    end = m.end(1) + offset
    
    new_html = new_html[:start] + new_path + new_html[end:]
    offset += len(new_path) - len(old_path)
    
    print(f"Replaced: {old_path} -> {new_path}")

with open(html_file, 'w', encoding='utf-8') as f:
    f.write(new_html)
    
print(f"Total images used from pool: {len(used)-1}")
print("Finished manual curation script.")
