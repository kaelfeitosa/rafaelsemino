import re

html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

# Replace all instances of cols-2 with cols-3 for carousels
html_content = html_content.replace('carousel-container--cols-2', 'carousel-container--cols-3')

with open(html_file, 'w', encoding='utf-8') as f:
    f.write(html_content)

print("Updated all carousels to use 3 columns.")
