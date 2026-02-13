import re
import os

def fix_previews(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Ensure all ![Preview](...) paths are relative and clean
    # Sometimes absolute paths or double slashes creep in.
    
    # Standardize to forward slashes
    content = content.replace('\\', '/')
    
    # Remove any file:/// prefixes that might have been added by accident in previous steps
    content = content.replace('file:///c:/Users/mkael/.gemini/antigravity/playground/warped-halley/', '')
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)

if __name__ == "__main__":
    selector_path = r'c:\Users\mkael\.\.gemini\antigravity\playground\warped-halley\image_selector.md'
    if os.path.exists(selector_path):
        fix_previews(selector_path)
        print("image_selector.md preview paths sanitized.")
    else:
        print(f"File not found: {selector_path}")
