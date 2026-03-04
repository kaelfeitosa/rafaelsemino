import os
import yaml
import re

acervo_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities"
html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
report_path = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\audit_report.txt"

with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

def get_frontmatter(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    if content.startswith('---'):
        parts = content.split('---')
        if len(parts) >= 3:
            try:
                fm = yaml.safe_load(parts[1])
                body = parts[2].strip()
                return fm, body
            except yaml.YAMLError:
                pass
    return None, None

with open(report_path, 'w', encoding='utf-8') as report:
    report.write("=========================================\n")
    report.write("ACERVO CONSISTENCY AUDIT REPORT\n")
    report.write("=========================================\n")

    for root, dirs, files in os.walk(acervo_dir):
        for filename in files:
            if filename.endswith('.md'):
                filepath = os.path.join(root, filename)
                meta, body = get_frontmatter(filepath)
                if not meta:
                    continue
                    
                draft = meta.get('draft', False)
                if draft:
                    continue
                    
                title = meta.get('title', '')
                name = meta.get('name', '')
                date = str(meta.get('date', ''))
                role = meta.get('role', '')
                
                search_term = title if title else name
                
                if search_term:
                    match = re.search(re.escape(search_term), html, re.IGNORECASE)
                    if match:
                        start = max(0, match.start() - 300)
                        end = min(len(html), match.end() + 300)
                        context = html[start:end]
                        
                        issues = []
                        if date and date not in context:
                            year_match = re.search(r'\b\d{4}\b', date)
                            if year_match:
                                year = year_match.group(0)
                                if year not in context:
                                    issues.append(f"Date/Year '{year}' missing/mismatch in HTML near this item.")
                            else:
                                issues.append(f"Date '{date}' missing/mismatch in HTML near this item.")
                                
                        if role:
                            role_keywords = role.split(',')
                            for r_k in role_keywords:
                                rk_clean = r_k.strip().split()[0] # first word
                                if rk_clean.lower() not in context.lower():
                                    issues.append(f"Role keyword '{rk_clean}' possibly missing in HTML near this item.")
                        
                        if issues:
                            report.write(f"\n[?] {filename}: {search_term}\n")
                            for i in issues:
                                report.write(f"    - {i}\n")
                    else:
                        report.write(f"\n[-] {filename}: {search_term}\n")
                        report.write(f"    - Item missing entirely from index.html\n")

    report.write("\n=========================================\n")
    report.write("Audit complete.\n")

print("Audit report written to audit_report.txt")
