import re

def main():
    with open('index.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Remove empty CSS ruleset
    # The review pointed to an empty ruleset.
    # Looking at the file content in previous turns, I see:
    # .clean-list,
    # .awards-list,
    # .eixos-list,
    # .direction-list {
    #     }
    # This was likely created when I removed list-style: none but kept the selector block accidentally empty?
    # Or maybe I moved properties to the shared block but left the original empty.
    # Wait, in refactor_final_v4.py I added a grouped rule for `list-style: none` AND removed it from individual blocks.
    # But I see in the `read_file` output:
    # .clean-list,
    # .awards-list,
    # .eixos-list,
    # .direction-list {
    # }
    # at the end of the style block? No, I see:
    # .clean-list li::before, ... { margin-right: 0.5em; }
    # .clean-list li::before { content: "- "; }
    # ...
    # And then `.catalog-grid--single-col`.
    # And then `.img-stack-item--fixed-height`.
    # And finally:
    # .clean-list,
    # .awards-list,
    # .eixos-list,
    # .direction-list {
    #     }
    # Yes, it is empty at the very end (from `refactor_final_v4.py` step 2).
    # Wait, in step 2 of `refactor_final_v4.py` I added:
    # grouped_list_style = """
    #     .clean-list,
    #     ... {
    #         list-style: none;
    #     }
    # """
    # So it shouldn't be empty.
    # Ah, maybe I messed up the regex replacement later?
    # Regardless, I need to remove this empty block.

    empty_ruleset_regex = r'\.clean-list,\s*\.awards-list,\s*\.eixos-list,\s*\.direction-list\s*\{\s*\}'
    content = re.sub(empty_ruleset_regex, '', content)

    # 2. Fix circular CSS variable definition
    # --text-light: var(--text-light); -> --text-light: #ddd;
    content = content.replace('--text-light: var(--text-light);', '--text-light: #ddd;')

    with open('index.html', 'w', encoding='utf-8') as f:
        f.write(content)

if __name__ == "__main__":
    main()
