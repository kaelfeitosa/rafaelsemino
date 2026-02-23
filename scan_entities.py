import os

ENTITIES_DIR = 'acervo/entities'

def parse_frontmatter(content):
    frontmatter = {}
    if content.startswith('---'):
        parts = content.split('---', 2)
        if len(parts) >= 3:
            lines = parts[1].strip().split('\n')
            for line in lines:
                if ':' in line:
                    key, value = line.split(':', 1)
                    frontmatter[key.strip()] = value.strip()
    return frontmatter

def scan_entities():
    for root, dirs, files in os.walk(ENTITIES_DIR):
        for file in files:
            if file.endswith('.md'):
                filepath = os.path.join(root, file)
                try:
                    with open(filepath, 'r') as f:
                        content = f.read()
                        frontmatter = parse_frontmatter(content)

                        folder_type = os.path.basename(root)
                        id_val = frontmatter.get('id', 'N/A')
                        type_val = frontmatter.get('type', 'N/A')
                        title = frontmatter.get('title', frontmatter.get('name', 'N/A'))
                        tags = frontmatter.get('tags', 'N/A')

                        print(f"File: {filepath}")
                        print(f"  Folder Type: {folder_type}")
                        print(f"  ID: {id_val}")
                        print(f"  Type: {type_val}")
                        print(f"  Title/Name: {title}")
                        print(f"  Tags: {tags}")
                        print("-" * 20)
                except Exception as e:
                    print(f"Error reading {filepath}: {e}")

if __name__ == "__main__":
    scan_entities()
