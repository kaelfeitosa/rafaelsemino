import os
import glob
from bs4 import BeautifulSoup
import re

html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
acervo_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities"

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

md_files = glob.glob(os.path.join(acervo_dir, '**/*.md'), recursive=True)

missing = []

for md_file in md_files:
    filename = os.path.basename(md_file).replace('.md', '')
    if 'template' in filename or 'index' in filename: continue
    
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Extract title
    title_match = re.search(r'^title:\s*(.*)', content, re.MULTILINE)
    title = title_match.group(1).strip().strip('"\'') if title_match else filename
    
    # Extract description or tags if any to judge "readiness"
    has_images = "[[work" in content or "[[agent" in content or "attach" in content.lower()
    content_length = len(content.splitlines())
    
    # Check if mentioned in HTML
    # We check the exact filename (ID) or the title
    if filename in html_content or title in html_content:
        pass
    else:
        # Check if parts of the title are in HTML (sometimes titles are rich, like "Exu Não Vem Hoje (2022)")
        title_stripped = title.split('(')[0].strip()
        if title_stripped and title_stripped in html_content:
            pass
        else:
            missing.append({
                'id': filename,
                'title': title,
                'has_images': has_images,
                'size': content_length,
                'path': md_file
            })

output_path = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\gap_results.txt"
with open(output_path, 'w', encoding='utf-8') as f:
    f.write(f"Total entities analyzed: {len(md_files)}\n")
    f.write(f"Missing entities found: {len(missing)}\n\n")

    f.write("--- POTENTIAL ADDITIONS TO SITE ---\n")
    for m in sorted(missing, key=lambda x: x['size'], reverse=True):
        img_status = "🖼 Has Images" if m['has_images'] else "📄 Text Only"
        f.write(f"📌 {m['title']}\n")
        f.write(f"   ID: {m['id']}\n")
        f.write(f"   Readiness: {img_status}, {m['size']} lines of markdown content\n")
        f.write(f"   Path: {m['path']}\n\n")

print(f"Results written to {output_path}")
