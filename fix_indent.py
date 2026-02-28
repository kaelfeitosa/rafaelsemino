import os
import glob

files = glob.glob('acervo/entities/**/*.md', recursive=True)
for f in files:
    with open(f, 'r') as file:
        lines = file.readlines()

    new_lines = []
    in_frontmatter = False
    count = 0
    changed = False
    for line in lines:
        if line.strip() == '---':
            count += 1
            if count == 1:
                in_frontmatter = True
            elif count == 2:
                in_frontmatter = False

        if in_frontmatter and line.startswith('    - caption:'):
            # this is overly indented, let's just make it 2 spaces
            line = '  - caption:' + line[14:]
            changed = True
        elif in_frontmatter and line.startswith('      role:'):
            line = '    role:' + line[11:]
            changed = True
        elif in_frontmatter and line.startswith('      src:'):
            line = '    src:' + line[10:]
            changed = True
        elif in_frontmatter and line.startswith('      type:'):
            line = '    type:' + line[11:]
            changed = True

        new_lines.append(line)

    if changed:
        with open(f, 'w') as file:
            file.writelines(new_lines)
