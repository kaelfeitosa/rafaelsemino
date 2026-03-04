import re
from bs4 import BeautifulSoup
import os

html_file = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_file, 'r', encoding='utf-8') as f:
    soup = BeautifulSoup(f.read(), 'html.parser')

cards = soup.find_all('div', class_='card')
borrowed = []

for card in cards:
    # Try to find the title of the card
    title_tag = card.find('h3')
    if title_tag:
        title = title_tag.get_text(strip=True)
    else:
        title_tag = card.find('h2')
        title = title_tag.get_text(strip=True) if title_tag else "Unknown Section"
    
    images = card.find_all('focus-image')
    for img in images:
        src = img.get('src')
        if src and src.startswith('images/optimized/'):
            folder = src.split('/')[2]
            
            # Heuristics to check if borrowed
            # For example, if we are in "Exu Não Vem Hoje", folder should be "work-exu-nao-vem-hoje"
            # It's hard to exactly match string to folder, so we just print them all out for review
            borrowed.append(f"Section: {title}\nImage: {folder}\nPath: {src}\n")

print("--- Image Audit ---")
for b in borrowed:
    print(b)
