import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# 1. Standardize casing in Atuação e Criação (Teatro Destaques)
html = html.replace('<strong>IRREVERSÍVEL</strong>', '<strong>Irreversível</strong>')
html = html.replace('<strong>SANTO BORDEL DE TIATIRA</strong>', '<strong>Santo Bordel de Tiatira</strong>')
html = html.replace('<strong>DE LOUCO, TODO MUNDO TEM UM POUCO</strong>', '<strong>De Louco, Todo Mundo Tem um Pouco</strong>')
html = html.replace('<strong>A SERPENTE</strong>', '<strong>A Serpente</strong>')
html = html.replace('<strong>DE SUCUPIRA À ASA BRANCA</strong>', '<strong>De Sucupira à Asa Branca</strong>')

# 2. Remove borrowed images from the "Atuação e Criação" consolidated card
# The card starts around <div class="card woodcut-border card--light">
# and has <div class="img-stack img-stack--cols-3">...</div>
img_stack_pattern = r'<div class="img-stack img-stack--cols-3">\s*<focus-image[^>]+work-exu-nao-vem-hoje[^>]+></focus-image>\s*<focus-image[^>]+work-exu-nao-vem-hoje[^>]+></focus-image>\s*<focus-image[^>]+work-exu-nao-vem-hoje[^>]+></focus-image>\s*</div>'
html = re.sub(img_stack_pattern, '', html, flags=re.DOTALL)

# 3. Move "Produção e Gestão" section to the end
prod_section_pattern = r'(<!-- SEÇÃO: PRODUÇÃO E GESTÃO -->\s*<section id="producao".*?</section>)'
prod_match = re.search(prod_section_pattern, html, flags=re.DOTALL)

if prod_match:
    prod_content = prod_match.group(1)
    # Remove from its current location
    html = html.replace(prod_content, '')
    
    # Ensure it's single-col
    prod_content = prod_content.replace('<div class="catalog-grid">', '<div class="catalog-grid catalog-grid--single-col gap-60">', 1)
    
    # Adjust the empty gap spacing if needed. Adding spacing.
    
    # Insert before </main>
    html = html.replace('</main>', prod_content + '\n\n        </main>')

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)

print("Layout adjustments applied.")
