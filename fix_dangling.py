import os
import re

entities_dirs = [
    "acervo/entities/actions",
    "acervo/entities/works",
    "acervo/entities/agents"
]

def get_entities():
    entities = {}
    for d in entities_dirs:
        if not os.path.exists(d): continue
        for f in os.listdir(d):
            if f.endswith('.md'):
                eid = f[:-3]
                entities[eid] = os.path.join(d, f)
    return entities

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

entities = get_entities()
entity_ids = set(entities.keys())
dangling = get_dangling_images()

print(f"Total dangling images: {len(dangling)}")

# Manual overrides
manual_map = {
    'agent-felipe-marques-001.jpeg': 'agent-coletivo-farol-novo',
    'agent-gabriel-franca-001.jpeg': 'agent-coletivo-farol-novo',
    'agent-zeis-001.jpeg': 'agent-coletivo-farol-novo',
    'agent-grupo-miraira-001.png': 'action-pesquisador-ifce',
    'agent-grupo-miraira-002.png': 'action-pesquisador-ifce',
    'agent-grupo-miraira-003.png': 'action-pesquisador-ifce',
    'agent-projeto-abarca-002.png': 'action-prof-percurso-basico',
    'work-blackheroes-001.jpeg': 'agent-rafael-semino',
    'work-blackheroes-002.jpeg': 'agent-rafael-semino',
    'work-blackheroes-003.png': 'agent-rafael-semino',
    'action-produtor-cultural-001.png': 'agent-rafael-semino',
    'action-apresentacao-centro-cultural-nordeste-001.jpeg': 'work-vao'
}

to_delete = [
    'action-rafael-teatro-2014-001.png',
    'action-rafael-teatro-2015-001.png'
]

skip = [
    'test-robust.jpeg'
]

file_updates = {}

for img in dangling:
    if img in skip:
        continue
    if img in to_delete:
        os.remove(os.path.join('acervo/media/images', img))
        print(f"DELETED: {img}")
        continue

    matched_eid = None
    if img in manual_map:
        matched_eid = manual_map[img]
    else:
        for eid in sorted(entity_ids, key=len, reverse=True):
            if img.startswith(eid):
                matched_eid = eid
                break

    if matched_eid:
        filepath = entities[matched_eid]
        if filepath not in file_updates:
            file_updates[filepath] = []
        file_updates[filepath].append(img)
        print(f"MATCH: {img} -> {matched_eid}")
    else:
        print(f"NO MATCH: {img}")

# Apply updates
for filepath, images in file_updates.items():
    with open(filepath, 'r') as f:
        content = f.read()

    # Simple YAML injection: Look for attachments block, or insert it.
    if 'attachments:' in content:
        # Find the attachments block and insert at the end of it
        lines = content.split('\n')
        insert_idx = -1
        in_attachments = False
        for i, line in enumerate(lines):
            if line.startswith('attachments:'):
                in_attachments = True
            elif in_attachments and not line.startswith('  ') and not line.startswith('-') and line.strip() != '' and not line.startswith('attachments:'):
                insert_idx = i
                break
            elif in_attachments and line == '---': # end of frontmatter
                insert_idx = i
                break

        if insert_idx != -1:
            for img in images:
                ext = img.split('.')[-1]
                mime = 'image'
                lines.insert(insert_idx, f"  - caption: {img}")
                lines.insert(insert_idx+1, f"    role: documentation")
                lines.insert(insert_idx+2, f"    src: {img}")
                lines.insert(insert_idx+3, f"    type: {mime}")
                insert_idx += 4
            content = '\n'.join(lines)
    else:
        # No attachments block, find --- and insert it before the second one
        lines = content.split('\n')
        # find second ---
        count = 0
        insert_idx = -1
        for i, line in enumerate(lines):
            if line == '---':
                count += 1
                if count == 2:
                    insert_idx = i
                    break
        if insert_idx != -1:
            lines.insert(insert_idx, "attachments:")
            insert_idx += 1
            for img in images:
                ext = img.split('.')[-1]
                mime = 'image'
                lines.insert(insert_idx, f"  - caption: {img}")
                lines.insert(insert_idx+1, f"    role: documentation")
                lines.insert(insert_idx+2, f"    src: {img}")
                lines.insert(insert_idx+3, f"    type: {mime}")
                insert_idx += 4
            content = '\n'.join(lines)

    with open(filepath, 'w') as f:
        f.write(content)

print("Done mapping images!")
