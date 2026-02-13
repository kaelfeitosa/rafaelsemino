import re
import os

def update_selector(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    new_lines = []
    
    # Pattern to match "- [ ] `images/path/to/img.ext`"
    # and "![Preview](images/path/to/img.ext)"
    
    i = 0
    while i < len(lines):
        line = lines[i]
        new_lines.append(line)
        
        # Check if line contains an image path
        match = re.search(r'- \[ \] `(images/.*?)`', line)
        if match:
            original_path = match.group(1)
            # Find the preview line which should be immediately after
            if i + 1 < len(lines) and '![Preview]' in lines[i+1]:
                preview_line = lines[i+1]
                new_lines.append(preview_line)
                
                # Construct cropped path
                # e.g. images/portfolio/img.jpg -> images/cropped/portfolio/img.jpg
                parts = original_path.split('/')
                # parts[0] is 'images'
                cropped_path = 'images/cropped/' + '/'.join(parts[1:])
                
                # Add the Alternative
                new_lines.append(f'- [ ] **Smart Crop** (Foco no detalhe): `{cropped_path}`\n')
                new_lines.append(f'![Preview]({cropped_path})\n')
                
                i += 1 # Skip the original preview line we already added
        i += 1
        
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(new_lines)

if __name__ == "__main__":
    selector_path = r'c:\Users\mkael\.\.gemini\antigravity\playground\warped-halley\image_selector.md'
    if os.path.exists(selector_path):
        update_selector(selector_path)
        print("image_selector.md updated with cropped alternatives.")
    else:
        print(f"File not found: {selector_path}")
