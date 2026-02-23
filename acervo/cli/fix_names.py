import subprocess
import os

os.chdir(r"c:\Users\mkael\progs\html\rafaelsemino\acervo\cli")

commands = [
    # Agents
    r'.\acervo.exe ingest update agent-ufc name="UFC - Universidade Federal do Ceará"',
    r'.\acervo.exe ingest update agent-ifce name="IFCE - Instituto Federal do Ceará"',
    r'.\acervo.exe ingest update agent-ufba name="UFBA - Universidade Federal da Bahia"',
    r'.\acervo.exe ingest update agent-ccbj name="CCBJ - Centro Cultural Bom Jardim"',
    r'.\acervo.exe ingest update agent-escola-paulo-petrola name="Escola Paulo Petrola"',
    r'.\acervo.exe ingest update agent-escola-hugo-sadrack name="Escola Hugo Sadrack do Vale"',

    # Events
    r'.\acervo.exe ingest update event-mestrado-ufc name="Mestrado em Artes"',
    r'.\acervo.exe ingest update event-licenciatura-ifce name="Licenciatura em Teatro"',
    r'.\acervo.exe ingest update event-pos-teatro-oprimido name="Pós-graduação em Teatro do Oprimido"',
    r'.\acervo.exe ingest update event-laboratorio-pesquisa-ccbj name="Laboratório de Pesquisa CCBJ"',
    r'.\acervo.exe ingest update event-aceleracao-idoso name="Programa de Aceleração do Idoso"',
    r'.\acervo.exe ingest update event-percurso-basico-teatro name="Percurso Básico de Teatro"'
]

for cmd in commands:
    subprocess.run(cmd, shell=True)

print("Nomes corrigidos.")
