import re

with open('acervo/entities/works/work-exu-nao-vem-hoje.md', 'r') as f:
    content = f.read()

# remove attachments that are missing
missing_images = [
    'work-exu-nao-vem-hoje-010.jpeg',
    'work-exu-nao-vem-hoje-003.jpeg',
    'work-exu-nao-vem-hoje-001.jpeg',
    'work-exu-nao-vem-hoje-005.png'
]

lines = content.split('\n')
new_lines = []
skip = False
for line in lines:
    if line.strip() == '- caption: Image':
        pass # we will check next lines

    if any(img in line for img in missing_images):
        if line.strip().startswith('src:'):
            # this means previous 2 lines were caption and role, pop them
            new_lines.pop()
            new_lines.pop()
            continue
        elif line.strip().startswith('!['):
            continue

    new_lines.append(line)

with open('acervo/entities/works/work-exu-nao-vem-hoje.md', 'w') as f:
    f.write('\n'.join(new_lines))
