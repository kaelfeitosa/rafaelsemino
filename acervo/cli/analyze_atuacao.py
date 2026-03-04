import os
import glob
import yaml

works_dir = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\entities\works"
report_path = r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli\atuacao_img_counts.txt"
files = glob.glob(os.path.join(works_dir, "*.md"))

lines = ["=== Atuação Projects ==="]
for f in files:
    with open(f, 'r', encoding='utf-8') as file:
        content = file.read()
    if content.startswith('---'):
        yaml_content = content.split('---')[1]
        try:
            data = yaml.safe_load(yaml_content)
            img_count = 0
            for key, val in data.items():
                if key.endswith('_type') and val == 'image':
                    img_count += 1
            
            role = str(data.get('role', '')).lower()
            medium = str(data.get('medium', '')).lower()
            if 'ator' in role or 'teatro' in medium or 'performance' in medium or 'atuacao' in medium:
                lines.append(f"{data.get('title', 'Unknown')} | Images: {img_count} | Role: {data.get('role')} | Medium: {data.get('medium')} | Year: {data.get('year')}")
        except:
            pass

with open(report_path, 'w', encoding='utf-8') as f:
    f.write("\n".join(lines))
