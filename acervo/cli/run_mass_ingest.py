import subprocess
import os

os.chdir(r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli")

commands = [
    # Cursos e Faculdades
    r'.\acervo.exe ingest create agent ufc name="UFC - Universidade Federal do Ceará"',
    r'.\acervo.exe ingest create event mestrado-ufc name="Mestrado em Artes" organizers="[[agent-ufc]]"',
    r'.\acervo.exe ingest create participation rafael-mestrado-ufc agent="[[agent-rafael-semino]]" event="[[event-mestrado-ufc]]" role="Mestrando"',

    r'.\acervo.exe ingest create agent ifce name="IFCE - Instituto Federal do Ceará"',
    r'.\acervo.exe ingest create event licenciatura-ifce name="Licenciatura em Teatro" organizers="[[agent-ifce]]"',
    r'.\acervo.exe ingest create participation rafael-licenciatura-ifce agent="[[agent-rafael-semino]]" event="[[event-licenciatura-ifce]]" role="Graduando"',

    r'.\acervo.exe ingest create agent ufba name="UFBA - Universidade Federal da Bahia"',
    r'.\acervo.exe ingest create event pos-teatro-oprimido name="Pós-graduação em Teatro do Oprimido" organizers="[[agent-ufba]]"',
    r'.\acervo.exe ingest create participation rafael-pos-ufba agent="[[agent-rafael-semino]]" event="[[event-pos-teatro-oprimido]]" role="Pós-graduando"',

    r'.\acervo.exe ingest create agent ccbj name="CCBJ - Centro Cultural Bom Jardim"',
    r'.\acervo.exe ingest create event laboratorio-pesquisa-ccbj name="Laboratório de Pesquisa CCBJ" organizers="[[agent-ccbj]]"',
    r'.\acervo.exe ingest create participation rafael-bolsa-ccbj agent="[[agent-rafael-semino]]" event="[[event-laboratorio-pesquisa-ccbj]]" role="Pesquisador"',

    # Pedagogia
    r'.\acervo.exe ingest create agent escola-paulo-petrola name="Escola Paulo Petrola"',
    r'.\acervo.exe ingest create participation prof-paulo-petrola agent="[[agent-rafael-semino]]" event="[[agent-escola-paulo-petrola]]" role="Professor de Artes"',

    r'.\acervo.exe ingest create agent escola-hugo-sadrack name="Escola Hugo Sadrack do Vale"',
    r'.\acervo.exe ingest create participation prof-hugo-sadrack agent="[[agent-rafael-semino]]" event="[[agent-escola-hugo-sadrack]]" role="Professor de Artes"',

    r'.\acervo.exe ingest create event aceleracao-idoso name="Programa de Aceleração do Idoso"',
    r'.\acervo.exe ingest create participation prof-aceleracao agent="[[agent-rafael-semino]]" event="[[event-aceleracao-idoso]]" role="Professor/Diretor"',

    r'.\acervo.exe ingest create event percurso-basico-teatro name="Percurso Básico de Teatro" organizers="[[agent-porto-iracema]]"',
    r'.\acervo.exe ingest create participation prof-percurso-basico agent="[[agent-rafael-semino]]" event="[[event-percurso-basico-teatro]]" role="Professor"',

    # Repertório (Still use title!)
    r'.\acervo.exe ingest create work irreversivel title="Irreversível" year="2022" language="teatro"',
    r'.\acervo.exe ingest create participation rafael-irreversivel agent="[[agent-rafael-semino]]" work="[[work-irreversivel]]" role="Ator"',

    r'.\acervo.exe ingest create work noite-de-alegria title="Noite de Alegria da Rua Trinta e Sete" language="teatro"',
    r'.\acervo.exe ingest create participation rafael-noite-alegria agent="[[agent-rafael-semino]]" work="[[work-noite-de-alegria]]" role="Intérprete/Orientador de Montagem"',

    r'.\acervo.exe ingest create work astronauta title="Astronauta" language="audiovisual"',
    r'.\acervo.exe ingest create participation rafael-astronauta agent="[[agent-rafael-semino]]" work="[[work-astronauta]]" role="Equipe (Angola)"'
]

print("Iniciando reconstrução purista via Schema no Acervo...")

for cmd in commands:
    print(f"Executando: {cmd}")
    try:
        subprocess.run(cmd, shell=True, check=True)
    except subprocess.CalledProcessError as e:
        print(f"Erro ao executar comando: {e}")
        break

print("Finalizado a reconstrução.")
