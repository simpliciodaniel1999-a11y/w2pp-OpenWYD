# Inventário de rotas e validação de destinos de teletransporte

**Estado:** investigação em andamento. Este documento registra evidências do código; não define ainda limites geométricos de plataformas nem autoriza integração no fluxo ativo.

## Fontes verificadas

- `Source/Code/TMSrv/GetFunc.cpp`, função `GetTeleportPosition`: tabela legada de origem, destino e condições.
- `Source/Code/TMSrv/_MSG_ReqTeleport.cpp`: cobra o custo e chama `DoTeleport` quando as coordenadas mudam.
- `tmserver/internal/world/teleport.go`: tabela atualmente portada para Go.
- `tmserver/internal/content/gamemap.go`: HeightMap de 4096×4096 e AttributeMap de 1024×1024; cada célula de atributo corresponde a um bloco de 4×4 células de altura.
- `tmserver/internal/route/route.go`: `route.Bake` aplica o bit 2 do AttributeMap e marca o bloco como `route.Blocked = 127`; `route.MH = 8` é uma regra de transição de altura do pathfinding, não um limite de plataforma.

## Como interpretar as coordenadas

A origem é identificada pelo tile arredondado para baixo ao múltiplo de 4 (`x & 0xFFFC`, `y & 0xFFFC`). O destino legado geralmente usa uma dispersão aleatória de 0 a 2 em cada eixo; algumas rotas usam 0 a 1. Por isso, validar somente a coordenada-base não é suficiente: é preciso examinar todas as células que a dispersão pode produzir.

**Importante:** estar dentro da grade 0–4095 e não ter o marcador `Blocked` não prova, isoladamente, que uma célula pertença à plataforma correta. A altura também participa da navegabilidade; além disso, uma célula livre próxima pode estar em outro lado de uma parede, borda ou divisão de mapa. Não usar o retângulo de uma cidade como limite genérico de dungeon, campo ou evento.

## Catálogo de rotas encontradas no código legado

Todas as origens e destinos abaixo são coordenadas do mundo (x, y). “Portada em Go” significa que a rota aparece na `teleportTable` atual, não que a segurança do destino já esteja validada.

| Origem (x,y) | Destino-base (x,y) | Rota/condição documentada no legado | Estado no Go |
|---|---:|---|---|
| (2116,2100) | (1044,1724) | Armia → Noatum; custo 700 | Portada |
| (2480,1716) | (1044,1716) | Azran → Noatum; custo 700 | Portada |
| (2456,2016) | (1044,1708) | Erion → Noatum; custo 700 | Portada |
| (1044,1724) | (2116,2100) | Noatum → Armia | Portada |
| (1044,1716) | (2480,1716) | Noatum → Azran | Portada |
| (1044,1708) | (2456,2016) | Noatum → Erion | Portada |
| (1048,1764) | (1100,1712) | Área de guilda de Noatum → área de guerra; verificar contexto de cerco | Ausente |
| (1052,1708) | (3650,3110) | Noatum → Karden/Nippleheim | Portada |
| (2140,2068) | (2588,2096) | Armia → campo de Armia | Portada |
| (2468,1716) | (2248,1556) | Azran → campo de Azran | Portada |
| (2364,2284) | (144,3788) | Campo de Armia → dungeon 1 | Portada |
| (144,3788) | (2364,2284) | Dungeon 1 → campo de Armia | Portada |
| (2668,2156) | (148,3774) | Campo de Armia → dungeon 1; dispersão de 0 a 1 | Portada |
| (144,3772) | (2668,2156) | Dungeon 1 → campo de Armia | Portada |
| (148,3780) | (1004,4028) | Dungeon 1 → dungeon 2 | Portada |
| (144,3780) | (1004,4028) | Dungeon 1 → dungeon 2 | Portada |
| (1004,4028) | (148,3780) | Dungeon 2 → dungeon 1 | Portada |
| (408,4072) | (1004,4064) | Dungeon 1 → dungeon 2 (rota alternativa) | Portada |
| (1004,4064) | (408,4072) | Dungeon 2 → dungeon 1 (rota alternativa) | Portada |
| (744,3820) | (1004,3992) | Dungeon 1 → dungeon 3 | Ausente |
| (1004,3992) | (744,3820) | Dungeon 3 → dungeon 1 | Ausente |
| (680,4076) | (916,3820) | Dungeon 2 → dungeon 3 | Ausente |
| (916,3820) | (680,4076) | Dungeon 3 → dungeon 2 | Ausente |
| (876,3872) | (932,3820) | Dungeon 2 → dungeon 3 | Ausente |
| (932,3820) | (876,3872) | Dungeon 3 → dungeon 2 | Ausente |
| (188,188) | (2548,1740) | Área de guerra → Azran | Ausente |
| (2548,1740) | (2281,3688) | Azran → Vale; exige item 3916 no slot 13 | Portada, condição parcial |
| (1824,1772) | (1172,4080) | Campo de Azran → submundo | Portada |
| (1172,4080) | (1824,1772) | Submundo → campo de Azran | Portada |
| (1516,3996) | (1304,3816) | Submundo, andar 1 → andar 2 | Ausente |
| (1304,3816) | (1516,3996) | Submundo, andar 2 → andar 1 | Ausente |
| (2452,1716) | (1969,1711) | Segundo teleporte de Azran → campo de Azran | Ausente |
| (2452,1988) | (1989,1755) | Teleporte de Erion → campo de Azran | Ausente |
| (3648,3108) | (1054,1710) | Karden/Nippleheim → Noatum | Portada |
| (1056,1724) | (1164,1720) | Noatum → RvR/deserto quando a guerra está inativa | Ausente; há destinos alternativos conforme estado e clã |
| (1056,1724) | (1061,2129) ou (1237,1966) | Noatum → RvR; guerra ativa, destino depende do clã | Ausente |
| (2364,3892) | (2364,3906) | Acesso ao Hall de Kefra; consome ticket | Ausente |
| (2364,3924) | (3250,1703) | Deserto → Kefra; depende de `KefraLive != 0` no legado | Portada, condição precisa ser comparada |
| (1312,1900) | (2366,4073) | Kefra → deserto | Ausente |
| (2364,4072) | (1314,1900) | Deserto → Kefra | Ausente |
| (3648,3108) | (1054,1710) | Karden → Noatum | Portada |

A lista é uma transcrição das ramificações identificadas na função consultada, não uma garantia de que todos os sistemas de transporte do cliente estejam cobertos. Algumas rotas têm destinos condicionais, rotas de item ou eventos que exigem ler também o estado do mundo antes de modelá-las em Go.

## O que já sabemos sobre limites

- A grade global possui dimensões fixas de 4096×4096.
- As zonas de cidade em `tmserver/internal/world/city.go` são retângulos de cidades (spawn/limite de cidade), não uma descrição completa de cada plataforma.
- `areaBox` e os retângulos do castelo são limites de eventos específicos e não devem ser reaproveitados como limites globais.
- `HeightMap.dat` e `AttributeMap.dat` estão versionados em `Release/TMsrv/run/`. Os arquivos são binários; para avaliar o terreno real é necessário lê-los com o carregador Go e aplicar o mesmo bake do servidor. O conector GitHub usado para esta inspeção não expõe o conteúdo binário como texto legível, então este documento não afirma valores de terreno que ainda não foram medidos.

## Plano de validação antes de integrar

1. Carregar os dois arquivos de mapa com `content.LoadGrid` e aplicar `route.Bake`, exatamente como na inicialização do servidor.
2. Para cada rota, analisar cada célula possível de dispersão do destino, registrando: coordenada, altura assinada, bloqueio pelo atributo, vizinhos alcançáveis segundo a regra de altura e ocupação dinâmica (quando houver mundo ativo).
3. Não inferir um retângulo de plataforma a partir de poucas células. Para uma região candidata, inspecionar o componente conectado de células caminháveis e comparar com os pontos de entrada/saída conhecidos. Confirmar regiões de dungeon/evento com dados ou regras do legado.
4. Adicionar testes com dados sintéticos para limites, regiões separadas por bloqueios e dispersões de destino; adicionar testes com mapa real apenas se os assets estiverem disponíveis no ambiente de CI de forma confiável.
5. Comparar rota a rota o comportamento legado e o Go, incluindo condições de item, ticket, estado do evento, clã e custo.
6. Só depois decidir se o resolvedor pode corrigir células ocupadas/bloqueadas, restringindo a busca à região validada da rota. Se não houver uma região validada ou nenhuma célula segura for encontrada, falhar de forma conservadora em vez de escolher uma célula arbitrária.

## Critério de conclusão

Não integrar `nearestTeleportCell` em `doTeleport` até que:
- o catálogo de rotas esteja reconciliado com o legado;
- os destinos e suas dispersões tenham sido verificados contra o mapa real;
- os limites regionais usados por cada rota tenham evidência documentada;
- testes de regressão cubram destinos válidos, bloqueados, ocupados, sem alternativa e condições de acesso.


## Problema funcional relatado: portais que não teleportam

Além da melhoria preventiva para evitar destinos inválidos, há um problema funcional relatado: **alguns portais atualmente não funcionam — o personagem entra/aciona o portal, mas não é teleportado**. Ainda não foram identificados aqui todos os portais afetados nem confirmada a causa.

Esse problema deve ser investigado separadamente da validação de células seguras. Uma rota ausente da tabela Go é uma hipótese a verificar, mas não prova por si só a causa: também é necessário conferir o reconhecimento da origem, as condições de acesso/evento, a execução do pedido de teleporte e a resposta enviada ao cliente.

### Procedimento de investigação

1. Registrar cada portal afetado com mapa/local, coordenadas aproximadas da origem, destino esperado e condições para ativá-lo.
2. Comparar a origem com `GetTeleportPosition` no legado e com `teleportTable`/`TeleportDestWithAccess` no Go.
3. Seguir o fluxo de `_MSG_ReqTeleport` até `doTeleport`, verificando se o pedido chega, se a rota é encontrada e se alguma condição interrompe a operação.
4. Adicionar testes de regressão para cada portal confirmado como defeituoso, sem alterar rotas que já funcionam.
5. Só depois de identificar a causa, corrigir a rota ou condição específica e validar que a correção não muda os destinos válidos existentes.

**Nota:** a lista de rotas marcadas como “Ausente” neste documento contém candidatas à comparação, não uma confirmação de que cada uma delas corresponde a um portal relatado como defeituoso.


### Relatos concretos para reprodução

- **Erion:** há vários portais na cidade; um deles não funciona corretamente. A coordenada exata e o destino esperado ainda precisam ser identificados. Não assumir que é o portal principal de Erion (2452,1988) → (1989,1755) sem reproduzir e confirmar no jogo.
- **Entrada de Kefra:** o acesso/portal de entrada não funciona corretamente. É necessário confirmar qual portal específico está sendo acionado e em que estado se encontra o evento de Kefra. O código legado registra um acesso ao Hall de Kefra em (2364,3892) → (2364,3906), com consumo de ticket, e outras rotas relacionadas ao deserto/Kefra; isso não basta para concluir que esse é o ponto exato relatado.

Prioridade de depuração: reproduzir os dois casos separadamente, registrar coordenada de origem, destino esperado, item/ticket necessário e estado do evento; então rastrear reconhecimento da rota, condições, execução de `doTeleport` e atualização enviada ao cliente. Não aplicar correção genérica nem inferir coordenadas sem confirmação.


## Diagnóstico após reproduzir os relatos

O relato do usuário esclarece que, tanto no portal de Erion quanto na entrada de Kefra, **não acontece nada**. A investigação no código encontrou um candidato concreto para Erion e adicionou instrumentação para identificar o caminho de Kefra:

- **Portal secundário de Erion:** o legado `GetTeleportPosition` contém a origem `(2452,1988)` com destino-base `(1989,1755)` e dispersão de `0..2` em cada eixo. Essa rota estava ausente de `tmserver/internal/world/teleport.go`. A rota foi portada e recebeu teste cobrindo todas as 16 posições dentro do tile de origem e a faixa de destinos.
- **Entrada do Hall de Kefra:** o legado só teleporta da origem `(2364,3892)` quando `KefraTicket > 0`, descontando uma entrada e enviando a contagem atualizada. O handler Go já tinha esse fluxo, mas retornava silenciosamente quando não havia entradas. Agora registra um aviso explícito nesse caso.
- **Portal do deserto para Kefra:** também depende de o estado do Kefra estar carregado e marcado como derrotado. O handler agora registra quando esse acesso é rejeitado por estado ausente/fechado.
- **Outras falhas de correspondência:** pedidos de teleporte que não correspondam a uma rota ou condição de acesso agora geram log com a posição atual, tile arredondado e dados relevantes de acesso. Isso deve revelar se o cliente está acionando o portal numa coordenada diferente da esperada.

A rota de Erion foi corrigida com base direta na tabela legada; ainda precisa passar pela CI e ser reproduzida no cliente. Para Kefra, os logs vão distinguir uma falta de entradas, estado do evento fechado e coordenada de origem que não corresponde ao portal. Não alteramos as regras de acesso para permitir entrada sem ticket nem com Kefra vivo.
