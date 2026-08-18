# Manual de Regras e Configurações - Damas Go

Este documento descreve detalhadamente as **regras do jogo de damas** adotadas no desenvolvimento do projeto, as **matrizes espaciais** disponíveis, o **sistema de notação algébrica** e as **configurações de jogo e inteligência artificial**.

---

## Sumário

- [1. Regras do Jogo Adotadas](#1-regras-do-jogo-adotadas)
  - [Conformidade com as Regras Brasileiras](#conformidade-com-as-regras-brasileiras)
  - [Movimentação das Peças Comuns (Pedras)](#movimentação-das-peças-comuns-pedras)
  - [Captura pelas Peças Comuns](#captura-pelas-peças-comuns)
  - [Lei da Maioria (Captura Obrigatória do Maior Número)](#lei-da-maioria-captura-obrigatória-do-maior-número)
  - [Capturas Múltiplas e em Cadeia](#capturas-múltiplas-e-em-cadeia)
  - [Promoção e Coroação a Dama](#promoção-e-coroação-a-dama)
  - [A Dama Voadora](#a-dama-voadora)
  - [Condições de Vitória e Empate](#condições-de-vitória-e-empate)
- [2. Matrizes Espaciais do Tabuleiro](#2-matrizes-espaciais-do-tabuleiro)
  - [Matriz 10x10 (Padrão Ampliado - Recomendado)](#matriz-10x10-padrão-ampliado---recomendado)
  - [Matriz 8x8 (Padrão Clássico)](#matriz-8x8-padrão-clássico)
  - [Comparativo das Matrizes](#comparativo-das-matrizes)
- [3. Notação Algébrica e Formatos de Entrada](#3-notação-algébrica-e-formatos-de-entrada)
  - [Sistema de Coordenadas Estilo Xadrez](#sistema-de-coordenadas-estilo-xadrez)
  - [Formatos de Entrada Aceitos](#formatos-de-entrada-aceitos)
  - [Exemplos Práticos](#exemplos-práticos)
- [4. Configurações de Jogo e Motores de IA](#4-configurações-de-jogo-e-motores-de-ia)
  - [Seleção dos Motores de Inteligência Artificial](#seleção-dos-motores-de-inteligência-artificial)
  - [Níveis de Dificuldade](#níveis-de-dificuldade)
  - [Simbologia e Guia Visual](#simbologia-e-guia-visual)

---

## 1. Regras do Jogo Adotadas

### Conformidade com as Regras Brasileiras
O jogo segue rigorosamente as **Regras Brasileiras de Damas**:

1. O jogo é praticado exclusivamente nas **casas escuras** (casas jogáveis do tabuleiro).
2. O jogador com as peças **Brancas** (`●` / `★`) sempre realiza o primeiro movimento da partida.
3. As peças **Pretas** (`○` / `☆`) são controladas pelo motor de Inteligência Artificial.

---

### Movimentação das Peças Comuns (Pedras)
* **Deslocamento Simples**: Uma pedra comum pode se mover apenas **1 casa na diagonal para frente**, desde que a casa de destino esteja livre.
* **Proibição de Recuo**: Em lances simples (sem captura), pedras comuns não podem recuar para trás.

---

### Captura pelas Peças Comuns
* **Direções de Captura**: A pedra comum pode capturar peças adversárias saltando **para frente E para trás** ao longo das diagonais.
* **Mecânica do Salto**: A peça salta sobre uma peça adversária adjacente e aterrissa imediatamente na casa diagonal seguinte, que deve estar vazia.
* **Remoção de Peças**: A peça capturada é retirada do tabuleiro.

---

### Lei da Maioria (Captura Obrigatória do Maior Número)
* A captura é **estritamente obrigatória**. Se houver possibilidade de captura, nenhum movimento simples é permitido.
* **Regra da Maioria**: Caso existam múltiplos caminhos de captura disponíveis para o jogador ou para a IA, é **obrigatório escolher a linha que capture a maior quantidade de peças adversárias**.
* O motor de regras do jogo valida automaticamente todas as possibilidades e impede lances que violem a Lei da Maioria, exibindo uma mensagem de aviso no painel de controle.

---

### Capturas Múltiplas e em Cadeia
* Quando uma peça atinge uma nova posição após um salto e dali for possível realizar outro salto de captura sobre outra peça adversária, a captura **deve continuar obrigatoriamente no mesmo lance** até que não haja mais peças ao alcance de captura.
* Durante uma combinação de capturas em cadeia, a mesma peça adversária não pode ser saltada mais de uma vez.

---

### Promoção e Coroação a Dama
* Quando uma pedra comum atinge a última fileira do tabuleiro do lado adversário:
  * **Brancas**: atingem a fileira superior (Linha 10 no tabuleiro 10x10 ou Linha 8 no 8x8).
  * **Pretas**: atingem a fileira inferior (Linha 1 no tabuleiro 10x10 ou 8x8).
* A peça é imediatamente promovida a **Dama** (`★` para Brancas, `☆` para Pretas), ganhando poderes especiais de movimentação e alcance.
* Se uma pedra atingir a linha de promoção durante uma captura em cadeia, ela só é coroada se encerrar o lance nessa linha.

---

### A Dama Voadora
Diferente das regras americanas (checkers), nas regras brasileiras a Dama é **voadora**:
* **Movimentação Livre**: Pode se deslocar por **qualquer número de casas vazias** ao longo de qualquer uma das 4 diagonais (para frente e para trás).
* **Captura à Distância**: Pode capturar uma peça adversária localizada a qualquer distância na diagonal, desde que não haja outras peças no caminho intermediário.
* **Escolha do Ponto de Parada**: Após saltar a peça adversária, a Dama pode pousar em **qualquer casa livre** imediatamente atrás ou várias casas adiante ao longo da mesma diagonal.

---

### Condições de Vitória e Empate

#### Vitória:
Um jogador vence a partida quando:
1. **Eliminação Total**: Capturar todas as peças do adversário.
2. **Bloqueio Total**: O adversário não possuir nenhum movimento legal disponível no seu turno (todas as peças restantes estão travadas).

#### Empate:
Uma partida termina empatada quando:
1. Atingir **40 meios-lances** (20 jogadas completas) consecutivos de Damas sem nenhuma captura de peça ou avanço de pedras comuns.

---

## 2. Matrizes Espaciais do Tabuleiro

O jogo oferece suporte a duas matrizes espaciais configuráveis:

### Matriz 10x10 (Padrão Ampliado - Recomendado)
* **Estrutura**: Grade com 10 colunas (**A** a **J**) e 10 linhas (**1** a **10**).
* **Casas**: 100 casas totais (50 jogáveis escuras e 50 claras).
* **Peças Iniciais**: 20 peças para cada jogador (4 fileiras iniciais completas).
  * Brancas: ocupam as linhas 1, 2, 3 e 4.
  * Pretas: ocupam as linhas 7, 8, 9 e 10.
  * Campo neutro inicial: linhas 5 e 6 livres.
* **Experiência**: Proporciona maior profundidade estratégica, maior número de linhas táticas e espaço ampliado para manobras de Damas voadoras.

---

### Matriz 8x8 (Padrão Clássico)
* **Estrutura**: Grade com 8 colunas (**A** a **H**) e 8 linhas (**1** a **8**).
* **Casas**: 64 casas totais (32 jogáveis escuras e 32 claras).
* **Peças Iniciais**: 12 peças para cada jogador (3 fileiras iniciais completas).
  * Brancas: ocupam as linhas 1, 2 e 3.
  * Pretas: ocupam as linhas 6, 7 e 8.
  * Campo neutro inicial: linhas 4 e 5 livres.
* **Experiência**: Partidas mais rápidas e com contato tático direto desde os primeiros turnos.

---

### Comparativo das Matrizes

| Característica | Matriz 10x10 (Ampliada) | Matriz 8x8 (Clássica) |
|---|:---:|:---:|
| **Dimensão Total** | 100 casas | 64 casas |
| **Casas Jogáveis** | 50 | 32 |
| **Peças por Jogador** | 20 | 12 |
| **Linhas Iniciais de Peças** | 4 fileiras | 3 fileiras |
| **Colunas** | A, B, C, D, E, F, G, H, I, J | A, B, C, D, E, F, G, H |
| **Linhas** | 1 a 10 | 1 a 8 |
| **Duração Média de Partida** | 30 a 60 lances | 15 a 35 lances |

---

## 3. Notação Algébrica e Formatos de Entrada

O jogo utiliza o **sistema de notação algébrica padrão de xadrez** para mapeamento das casas, permitindo que o jogador informe seus lances de forma natural e intuitiva.

```
       A   B   C   D   E   F   G   H   I   J
  10  [ ] [○] [ ] [○] [ ] [○] [ ] [○] [ ] [○]  10
   9  [○] [ ] [○] [ ] [○] [ ] [○] [ ] [○] [ ]   9
   8  [ ] [○] [ ] [○] [ ] [○] [ ] [○] [ ] [○]   8
   7  [○] [ ] [○] [ ] [○] [ ] [○] [ ] [○] [ ]   7
   6  [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]   6
   5  [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]   5
   4  [ ] [●] [ ] [●] [ ] [●] [ ] [●] [ ] [●]   4
   3  [●] [ ] [●] [ ] [●] [ ] [●] [ ] [●] [ ]   3
   2  [ ] [●] [ ] [●] [ ] [●] [ ] [●] [ ] [●]   2
   1  [●] [ ] [●] [ ] [●] [ ] [●] [ ] [●] [ ]   1
       A   B   C   D   E   F   G   H   I   J
```

---

### Formatos de Entrada Aceitos

O parser inteligente ([`pkg/game/notation.go`](file:///pkg/game/notation.go)) aceita qualquer um dos seguintes formatos (independente de maiúsculas ou minúsculas):

1. **Por Extenso em Português / Inglês**:
   * `E3 para F4` ou `e3 para f4`
   * `E3 to F4` ou `e3 to f4`
2. **Direto com Espaço ou Hífen**:
   * `E3 F4` ou `e3 f4`
   * `E3-F4` ou `E3->F4`
3. **Notação de Captura Simples**:
   * `E3 x G5` ou `E3:G5`
   * `E3 para G5` ou `E3 G5`
4. **Notação de Capturas em Cadeia**:
   * `C3:E5:G7` ou `C3 para E5 para G7`
   * `C3 G7` (o motor encontra automaticamente o caminho de salto completo)

---

### Exemplos Práticos

| Situação | Digitação no Terminal | Ação Executada |
|---|---|---|
| Movimento Simples | `C3 para D4` ou `C3 D4` | Avança a pedra de C3 para D4 |
| Captura Simples | `E3 para G5` ou `E3:G5` | Salta sobre a peça em F4 e aterrissa em G5 |
| Salto Duplo em Cadeia | `C3:E5:G7` ou `C3 G7` | Captura a peça em D4 (pousa em E5) e em F6 (pousa em G7) |
| Lance com Dama Voadora | `A1 para H8` ou `A1 H8` | Voa pela diagonal principal de A1 até H8 |
| Encerrar Partida | `sair`, `exit` ou `q` | Encerra o jogo e retorna ao terminal |

---

## 4. Configurações de Jogo e Motores de IA

Ao iniciar o jogo, o usuário é apresentado a menus interativos com opções de personalização:

### Seleção dos Motores de Inteligência Artificial

```
┌────────────────────────────────────────────────────────────────────────┐
│                      MOTOR DE INTELIGENCIA ARTIFICIAL                  │
├───────┬───────────────────────────────────┬────────────────────────────┤
│ OPCAO │ MOTOR DE IA                       │ CARACTERISTICA             │
├───────┼───────────────────────────────────┼────────────────────────────┤
│ 1     │ Modo Hibrido Mestre [Recomendado] │ Combina A*, MDP e HC       │
│ 2     │ Processo de Decisao de Markov     │ Value Iteration & Softmax  │
│ 3     │ Busca A* (A-Star)                 │ Táticas com Min-Heap       │
│ 4     │ Hill Climbing com Reinicializacao │ Otimizacao & Restarts      │
└───────┴───────────────────────────────────┴────────────────────────────┘
```

1. **Modo Híbrido Mestre**: Seleciona a estratégia adequada para cada fase (A* para táticas, MDP para postura de jogo, Hill Climbing para desempate).
2. **Processo de Decisão de Markov (MDP)**: Avalia probabilidades estocásticas de resposta do jogador humano com desconto temporal $\gamma = 0.90$.
3. **Busca A\***: Prioriza variantes de ganho forçado e cálculo de profundidade tática.
4. **Hill Climbing com Reinicialização Aleatória**: Realiza busca local heurística explorando planos e executando múltiplos restarts contra máximos locais.

---

### Níveis de Dificuldade

O nível de dificuldade ajusta dinamicamente a profundidade e a intensidade computacional dos algoritmos:

```
┌────────────────────────────────────────────────────────────────────────┐
│                          NIVEL DE DIFICULDADE                          │
├───────┬──────────┬─────────────────────────────────────────────────────┤
│ OPCAO │ NIVEL    │ PARAMETROS DE COMPUTACAO                            │
├───────┼──────────┼─────────────────────────────────────────────────────┤
│ 1     │ Fácil    │ Profundidade = 2 | Nós A* = 250 | Restarts HC = 10  │
│ 2     │ Médio    │ Profundidade = 3 | Nós A* = 600 | Restarts HC = 20  │
│ 3     │ Difícil  │ Profundidade = 4 | Nós A* = 1200 | Restarts HC = 40 │
└───────┴──────────┴─────────────────────────────────────────────────────┘
```

---

### Simbologia e Guia Visual

| Elemento | Símbolo | Cor no Terminal | Descrição |
|---|:---:|---|---|
| **Pedra Branca** | `●` | Ciano Brilhante | Peça comum do jogador |
| **Dama Branca** | `★` | Amarelo Brilhante | Dama coroada do jogador |
| **Pedra Preta** | `○` | Vermelho Brilhante | Peça comum da IA |
| **Dama Preta** | `☆` | Magenta Brilhante | Dama coroada da IA |
| **Casa Escura** | `   ` | Fundo Cinza Escuro | Casa jogável |
| **Casa Clara** | `   ` | Fundo Cinza Claro | Casa não-jogável |
| **Último Lance** | `   ` | Fundo Dourado/Oliva | Destaca a origem e destino da jogada anterior |
