# Puzzle_66_Solver

## Descrição

Em 2015, para demonstrar a imensidão do espaço da chave privada (ou talvez apenas por diversão), alguém criou um “quebra-cabeça/puzzle” onde escolhia chaves em um determinado espaço menor e enviava quantidades crescentes de bitcoin para cada uma dessas chaves.

Para resolver o quebra-cabeça, você deve iterar no espaço específico da chave privada e verificar o saldo de cada chave. Quanto mais estreito for o espaço da chave, maior será a chance de encontrar a chave correta. Atualmente, é melhor focar na resolução do quebra-cabeça nº 66 devido ao seu espaço de chave mais estreito.

## Objetivo

Este script foca na resolução do Puzzle 66. A chave privada para este puzzle está no intervalo de `20000000000000000` até `3ffffffffffffffff`. Com a chave correta, é possível obter a recompensa do puzzle.

## Funcionamento

- O script utiliza a CPU para maximizar o desempenho, utilizando multiprocessamento e todos os núcleos do processador disponíveis.
- Soluções com GPU são mais eficientes, mas também mais caras.
- O script gera uma chave privada hexadecimal aleatória dentro do intervalo especificado e realiza uma varredura de 10 milhões de chaves.
- Ao finalizar a varredura, o processo é repetido indefinidamente até encontrar a chave correta.
- O script tem como alvo a chave pública no formato `hash160`, evitando a necessidade de usar `base58` para obter o endereço `p2pkh`, priorizando a velocidade e evitando conversões criptográficas desnecessárias.

