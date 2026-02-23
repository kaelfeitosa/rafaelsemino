import os
import re

ENTITIES_DIR = 'acervo/entities'

# Regex for parsing markdown links like [[link]]
LINK_REGEX = re.compile(r'\[\[(.*?)\]\]')

def parse_frontmatter(content):
    frontmatter = {}
    if content.startswith('---'):
        parts = content.split('---', 2)
        if len(parts) >= 3:
            lines = parts[1].strip().split('\n')
            for line in lines:
                if ':' in line:
                    key, value = line.split(':', 1)
                    key = key.strip()
                    value = value.strip()
                    # Handle basic list parsing [val1, val2]
                    if value.startswith('[') and value.endswith(']'):
                        # Remove brackets and quotes
                        items = value[1:-1].split(',')
                        clean_items = []
                        for item in items:
                            item = item.strip().strip('"\'')
                            if item:
                                clean_items.append(item)
                        value = clean_items
                    frontmatter[key] = value
    return frontmatter

def get_linked_id(val):
    if isinstance(val, list):
        return [get_linked_id(v) for v in val]
    if isinstance(val, str):
        match = LINK_REGEX.search(val)
        if match:
            return match.group(1)
        return val
    return val

def scan_entities():
    entities = {}

    # 1. Load all entities
    for root, dirs, files in os.walk(ENTITIES_DIR):
        for file in files:
            if file.endswith('.md'):
                filepath = os.path.join(root, file)
                try:
                    with open(filepath, 'r') as f:
                        content = f.read()
                        fm = parse_frontmatter(content)
                        entity_id = fm.get('id', os.path.splitext(file)[0])
                        # If ID is missing in frontmatter, use filename
                        if not entity_id:
                            entity_id = os.path.splitext(file)[0]

                        entities[entity_id] = {
                            'id': entity_id,
                            'type': fm.get('type', 'unknown'),
                            'data': fm,
                            'filepath': filepath
                        }
                except Exception as e:
                    print(f"Error reading {filepath}: {e}")

    issues = []

    # 2. Check rules
    for eid, ent in entities.items():
        etype = ent['type']
        data = ent['data']

        # Mandatory Fields
        missing = []
        if etype == 'agent':
            if 'name' not in data and 'title' not in data: missing.append('name/title')
        elif etype == 'work':
            if 'title' not in data: missing.append('title')
            # created_by is mandatory but might be missing in some test data
            if 'created_by' not in data: missing.append('created_by')
        elif etype == 'event':
            if 'title' not in data and 'name' not in data: missing.append('title/name')
            if 'date_start' not in data and 'date' not in data: missing.append('date_start/date')
            if 'location' not in data: missing.append('location')
        elif etype == 'participation':
            if 'agent' not in data: missing.append('agent')
            # Participation must link to Event OR Work (sometimes directly to work?)
            # Audit rules say: Participation takes_place_in Event. But some might link to Work.
            if 'event' not in data and 'work' not in data: missing.append('event/work')
            if 'role' not in data: missing.append('role')
        elif etype == 'record':
            if 'related_to' not in data: missing.append('related_to')

        if missing:
            issues.append(f"[MISSING FIELDS] Entity {eid} ({etype}) missing: {', '.join(missing)}")

        # Invalid Relations (Direct links)
        # Agent -> Event direct?
        if etype == 'agent' and ('event' in data or 'events' in data):
            issues.append(f"[FORBIDDEN REL] Entity {eid} (agent) links directly to Event.")

        # Work -> Event direct?
        if etype == 'work' and ('event' in data or 'events' in data):
             issues.append(f"[FORBIDDEN REL] Entity {eid} (work) links directly to Event.")

        # Record -> Agent direct?
        if etype == 'record':
            related = data.get('related_to', [])
            if isinstance(related, str): related = [related]
            for rel in related:
                rel_id = get_linked_id(rel)
                # Check if rel_id exists and is an Agent
                if rel_id in entities:
                    target_type = entities[rel_id]['type']
                    if target_type == 'agent':
                        issues.append(f"[SUSPICIOUS RECORD] Entity {eid} (record) links directly to Agent {rel_id}. Should link to Participation.")
                elif rel_id:
                    issues.append(f"[BROKEN LINK] Entity {eid} (record) links to non-existent {rel_id}.")

    # 3. Duplicate Participations
    # Key: (AgentID, EventID, Role) -> [ParticipationIDs]
    participations_map = {}
    for eid, ent in entities.items():
        if ent['type'] == 'participation':
            data = ent['data']
            agent = get_linked_id(data.get('agent', ''))
            event = get_linked_id(data.get('event', ''))
            # Sometimes 'event' is a list? Assuming string for now based on previous files
            if isinstance(agent, list) and agent: agent = agent[0]
            if isinstance(event, list) and event: event = event[0]

            role = data.get('role', '')

            if agent and event and role:
                key = (agent, event, role)
                if key not in participations_map:
                    participations_map[key] = []
                participations_map[key].append(eid)
            elif agent and not event:
                # Check if linked to work instead?
                work = get_linked_id(data.get('work', ''))
                if isinstance(work, list) and work: work = work[0]
                if work:
                    key = (agent, work, role) # Different key for work-based participation
                    if key not in participations_map:
                        participations_map[key] = []
                    participations_map[key].append(eid)

    for key, p_ids in participations_map.items():
        if len(p_ids) > 1:
            issues.append(f"[DUPLICATE PARTICIPATION] (Agent: {key[0]}, Target: {key[1]}, Role: {key[2]}) found in: {', '.join(p_ids)}")

    # Print Report
    print("=== AUDIT REPORT ===")
    if not issues:
        print("No structural issues found.")
    else:
        for issue in issues:
            print(issue)

if __name__ == "__main__":
    scan_entities()
