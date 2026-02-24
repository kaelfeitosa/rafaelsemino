import os
import yaml
import json
import collections

# Config
ENTITIES_DIR = "acervo/entities"
MEDIA_DIR = "acervo/media/images"

# Report structure
report = {
    "structures": {"errors": [], "warnings": []},
    "actions": {"errors": [], "warnings": []},
    "evidence": {"errors": [], "warnings": []},
    "narrative": {"timeline": [], "gaps": []}
}

# Counters
stats = collections.defaultdict(int)

# Global entity store
entities = {}

def load_entities():
    """Reads all YAML files in entities/ subdirs."""
    for root, _, files in os.walk(ENTITIES_DIR):
        for file in files:
            if not file.endswith(".md"):
                continue
            path = os.path.join(root, file)
            try:
                with open(path, "r", encoding="utf-8") as f:
                    content = f.read()
                    # Extract frontmatter
                    if content.startswith("---"):
                        parts = content.split("---", 2)
                        if len(parts) >= 3:
                            fm = yaml.safe_load(parts[1])
                            body = parts[2]

                            if not fm:
                                report["structures"]["errors"].append(f"Empty frontmatter: {file}")
                                continue

                            eid = fm.get("id")
                            if not eid:
                                report["structures"]["errors"].append(f"Missing ID: {file}")
                                continue

                            # Determine type based on folder
                            folder = os.path.basename(root)
                            # Or infer from content if needed

                            if eid in entities:
                                report["structures"]["errors"].append(f"Duplicate ID: {eid} in {file} and {entities[eid]['path']}")
                            else:
                                entities[eid] = {
                                    "path": path,
                                    "fm": fm,
                                    "body": body,
                                    "folder": folder
                                }
                                stats[folder] += 1
                        else:
                             report["structures"]["errors"].append(f"Invalid YAML block structure: {file}")
                    else:
                        report["structures"]["errors"].append(f"No YAML frontmatter: {file}")
            except Exception as e:
                report["structures"]["errors"].append(f"Error reading {file}: {str(e)}")

def check_structures():
    """Layer 1: Structural Validity"""
    for eid, data in entities.items():
        fm = data["fm"]
        folder = data["folder"]

        # Check mandatory fields per type
        if folder == "agents":
            required = ["id", "name", "kind"]
            for f in required:
                if f not in fm:
                    report["structures"]["errors"].append(f"Agent {eid} missing '{f}'")
            if fm.get("kind") not in ["person", "collective"]:
                 report["structures"]["errors"].append(f"Agent {eid} invalid kind: {fm.get('kind')}")

        elif folder == "works":
            required = ["id", "title", "type"]
            for f in required:
                if f not in fm:
                    report["structures"]["errors"].append(f"Work {eid} missing '{f}'")

        elif folder == "actions":
            required = ["id", "title", "performed_by", "my_role", "context", "date_start"]
            for f in required:
                if f not in fm:
                    report["structures"]["errors"].append(f"Action {eid} missing '{f}'")

            # Check links
            pb = fm.get("performed_by", "")
            if pb:
                # remove [[ ]]
                pb_id = pb.strip("[]")
                # Handle cases where [[ ]] might be missing or different
                if pb.startswith("[[") and pb.endswith("]]"):
                    pb_id = pb[2:-2]

                if pb_id not in entities:
                     report["structures"]["errors"].append(f"Action {eid} performed_by unknown agent: {pb} (parsed as {pb_id})")

            wid = fm.get("work_id", "")
            if wid:
                 w_id = wid
                 if wid.startswith("[[") and wid.endswith("]]"):
                     w_id = wid[2:-2]

                 if w_id not in entities:
                      report["structures"]["errors"].append(f"Action {eid} refers to unknown work: {wid}")

def check_actions():
    """Layer 2: Action Semantics"""
    titles = collections.defaultdict(list)

    for eid, data in entities.items():
        if data["folder"] != "actions":
            continue

        fm = data["fm"]

        # Check for ambiguity
        if not fm.get("my_role"):
             report["actions"]["errors"].append(f"Action {eid} missing 'my_role' (critical for attribution)")

        # Check context
        ctx = fm.get("context")
        if not ctx or not isinstance(ctx, dict) or not ctx.get("label"):
             report["actions"]["errors"].append(f"Action {eid} has invalid/missing context")

        # Duplicate detection (fuzzy)
        title = fm.get("title", "").lower().strip()
        titles[title].append(eid)

    for title, eids in titles.items():
        if len(eids) > 1:
            report["actions"]["warnings"].append(f"Potential duplicate action title '{title}': {eids}")

def check_evidence():
    """Layer 3: Evidence (Attachments)"""
    # Build list of all referenced images to check existence
    referenced_images = set()

    for eid, data in entities.items():
        fm = data["fm"]
        atts = fm.get("attachments", [])

        if data["folder"] in ["works", "actions"]:
            if not atts:
                report["evidence"]["warnings"].append(f"{data['folder'].capitalize()} {eid} has no attachments")
            else:
                for att in atts:
                    src = att.get("src")
                    if src:
                        full_path = os.path.join(MEDIA_DIR, src)
                        if not os.path.exists(full_path):
                            report["evidence"]["errors"].append(f"{eid} references missing image: {src}")
                        referenced_images.add(src)

                        # Check type vs placement? Hard to automate "malpositioned" without AI

    # Check for redundant usage (same image in multiple entities)
    # This might be okay, but worth flagging if excessive
    image_usage = collections.defaultdict(list)
    for eid, data in entities.items():
        atts = data["fm"].get("attachments", [])
        for att in atts:
             src = att.get("src")
             if src:
                 image_usage[src].append(eid)

    for img, eids in image_usage.items():
        if len(eids) > 1:
            # If used in more than 2 places, warn
             report["evidence"]["warnings"].append(f"Image {img} used in multiple entities: {eids}")

def check_narrative():
    """Layer 4: Narrative Timeline"""
    timeline = []
    for eid, data in entities.items():
        fm = data["fm"]
        date = fm.get("date_start") or fm.get("year")
        if date:
            timeline.append({
                "date": str(date),
                "id": eid,
                "title": fm.get("title") or fm.get("name"),
                "type": data["folder"]
            })

    # Sort
    timeline.sort(key=lambda x: x["date"])
    report["narrative"]["timeline"] = timeline

    # Simple gap analysis? (e.g. year jumps > 2)
    last_year = None
    for item in timeline:
        try:
            # simple assumption: date starts with YYYY
            year = int(item["date"][:4])
            if last_year and (year - last_year > 2):
                 report["narrative"]["gaps"].append(f"Gap between {last_year} and {year}")
            last_year = year
        except:
            pass

def main():
    load_entities()
    check_structures()
    check_actions()
    check_evidence()
    check_narrative()

    print(json.dumps(report, indent=2, ensure_ascii=False))

    with open("audit_results.json", "w", encoding="utf-8") as f:
        json.dump(report, f, indent=2, ensure_ascii=False)

if __name__ == "__main__":
    main()
