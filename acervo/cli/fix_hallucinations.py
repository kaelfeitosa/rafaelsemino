import re

html_path = r"c:\Users\mkael\progs\html\rafaelsemino\frontend\index.html"
with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

# Fix ABARCA / Hugo Sadrack Visual Card
target_1 = """                        <!-- ABARCA -->
                        <div class="card woodcut-border card--light">
                            <div class="img-stack img-stack--cols-2">
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-001.webp"
                                    alt="Aula do Projeto ABARCA"></focus-image>
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-002.webp"
                                    alt="Atividade prática com alunos"></focus-image>
                            </div>
                            <div class="card-content">
                                <span class="card-tag card-tag--red">Escola Porto Iracema das Artes 📍 Fortaleza – CE –
                                    Brasil</span>
                                <h3 class="course-institution">Projeto ABARCA</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Territórios de atuação:</strong> Vicente Pinzón,
                                        Genibaú, Itapipoca.</p>
                                    <p class="course-details mt-20"><strong>Disciplinas ministradas:</strong> Teatro,
                                        Dramaturgia, Construção de personagem, Jogos teatrais, Processos de criação.</p>
                                </div>
                            </div>
                        </div>"""

replacement_1 = """                        <!-- HUGO SADRACK -->
                        <div class="card woodcut-border card--light">
                            <div class="img-stack img-stack--cols-2">
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-001.webp"
                                    alt="Estudantes em ação"></focus-image>
                                <focus-image
                                    src="images/optimized/work-prof-hugo-sadrack/work-prof-hugo-sadrack-002.webp"
                                    alt="Atividade prática com alunos"></focus-image>
                            </div>
                            <div class="card-content">
                                <span class="card-tag card-tag--red">Educação Básica (2021)</span>
                                <h3 class="course-institution">Escola Estadual Hugo Sadrack do Vale</h3>
                                <div class="mt-20">
                                    <p class="course-details"><strong>Disciplinas ministradas:</strong> Artes, Jogos e Africanidade, e História da Arte.</p>
                                    <p class="course-details mt-20">Atuação para turmas de Ensino Fundamental e Ensino Médio.</p>
                                </div>
                            </div>
                        </div>"""
                        
html = html.replace(target_1, replacement_1)

# Fix ABARCA / Hugo Sadrack Text List
target_2 = """                                    <!-- EDUCAÇÃO BÁSICA -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Educação Básica</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor de Artes</strong> – Escola Paulo Petrola</li>
                                            <li><strong>Professor de Artes</strong> – Escola Hugo Sadrack do Vale</li>
                                        </ul>
                                    </div>

                                    <!-- MEDIAÇÃO E CURSOS LIVRES -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Mediação e Cursos Livres</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor (Percurso Básico de Teatro)</strong> – Porto Iracema
                                                das Artes</li>"""
                                                
replacement_2 = """                                    <!-- EDUCAÇÃO BÁSICA -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Educação Básica</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor de Artes</strong> – Escola Municipal Paulo Petrola</li>
                                        </ul>
                                    </div>

                                    <!-- MEDIAÇÃO E CURSOS LIVRES -->
                                    <div class="atuacao-group">
                                        <p class="group-label">Mediação e Cursos Livres</p>
                                        <ul class="detail-list">
                                            <li><strong>Professor (Projeto ABARCA e Percurso Básico de Teatro)</strong> – Porto Iracema
                                                das Artes</li>"""
                                                
html = html.replace(target_2, replacement_2)

# Fix Azusa and Black Heroes
target_3 = """                            <div class="atuacao-consolidated-grid">
                                <!-- COORDENAÇÃO DE FORMAÇÃO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Coordenação de Formação</p>
                                    <ul class="detail-list">
                                        <li><strong>Black Heroes (2018–2023)</strong> – Produtor Cultural e Coordenador
                                            de Formação</li>
                                        <li><strong>Azusa (2018–2023)</strong> – Coordenação de ensino, pesquisa e
                                            processos formativos</li>
                                    </ul>
                                </div>
                            </div>"""
                            
replacement_3 = """                            <div class="atuacao-consolidated-grid">
                                <!-- PRODUÇÃO CULTURAL E COORDENAÇÃO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Produção Cultural e Coordenação</p>
                                    <ul class="detail-list">
                                        <li><strong>Black Heroes (2023)</strong> – Produtor Cultural e Coordenador de Formação</li>
                                        <li><strong>Rua Azusa (2018)</strong> – Produtor de espetáculo de teatro musical</li>
                                    </ul>
                                </div>
                            </div>"""

html = html.replace(target_3, replacement_3)

# Additional pass on Black Heroes/Azusa just in case there were formatting variations
target_3_alt = """                            <div class="atuacao-consolidated-grid">
                                <!-- COORDENAÇÃO DE FORMAÇÃO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Coordenação de Formação</p>
                                    <ul class="detail-list">
                                        <li><strong>Black Heroes (2018–2023)</strong> – Produtor Cultural e Coordenador
                                            de Formação</li>
                                        <li><strong>Azusa (2018–2023)</strong> – Coordenação de ensino, pesquisa e
                                            processos formativos</li>
                                    </ul>
                                </div>
                            </div>"""

html = re.sub(r'<!-- COORDENAÇÃO DE FORMAÇÃO -->.*Azusa \(2018–2023.*</ul>\s*</div>', 
              r'''<!-- PRODUÇÃO CULTURAL E COORDENAÇÃO -->
                                <div class="atuacao-group">
                                    <p class="group-label">Produção Cultural e Coordenação</p>
                                    <ul class="detail-list">
                                        <li><strong>Black Heroes (2023)</strong> – Produtor Cultural e Coordenador de Formação</li>
                                        <li><strong>Rua Azusa (2018)</strong> – Produtor de espetáculo de teatro musical</li>
                                    </ul>
                                </div>''', html, flags=re.DOTALL)

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)
print("Changes applied!")
