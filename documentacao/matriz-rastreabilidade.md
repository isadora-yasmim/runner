# Matriz de Rastreabilidade

Este documento relaciona os requisitos funcionais do Sistema Runner com os artefatos de planejamento e desenvolvimento do projeto, a fim de facilitar o acompanhamento da implementação, dos testes e da documentação.

## Requisitos funcionais e rastreamento

| Requisito | Descrição | Épico relacionado | Artefatos relacionados | Situação |
|---|---|---|---|---|
| US-01 | Invocar `assinador.jar` via CLI | Épico 1 — Assinador<br>Épico 2 — Servidor HTTP do Assinador<br>Épico 3 — CLI Assinatura | `documentacao/especificacao.md`<br>`documentacao/design.md`<br>`gerenciamento/backlog.md` | Planejado |
| US-02 | Simular assinatura digital com validação de parâmetros | Épico 1 — Assinador | `documentacao/especificacao.md`<br>`documentacao/definicao.md`<br>`gerenciamento/backlog.md` | Planejado |
| US-03 | Gerenciar ciclo de vida do Simulador do HubSaúde | Épico 4 — Simulador HubSaúde | `documentacao/especificacao.md`<br>`documentacao/design.md`<br>`gerenciamento/backlog.md` | Planejado |
| US-04 | Provisionar JDK automaticamente | Épico 5 — Provisionamento automático do JDK | `documentacao/especificacao.md`<br>`gerenciamento/backlog.md` | Planejado |
| US-05 | Disponibilizar binários multiplataforma | Épico 6 — Qualidade, Testes e Documentação | `documentacao/especificacao.md`<br>`gerenciamento/backlog.md`<br>`cronograma de execucao.md` | Planejado |

## Critérios de aceitação e cobertura prevista

| Requisito | Critérios de aceitação | Cobertura prevista |
|---|---|---|
| US-01 | Comandos de criação e validação via CLI; invocação direta e via HTTP; exibição legível do resultado | Testes de integração e testes de aceitação |
| US-02 | Validação rigorosa de parâmetros; simulação de criação e validação; mensagens claras de erro | Testes unitários, testes de integração e cenários de erro |
| US-03 | Iniciar, parar e monitorar o simulador; verificar portas; baixar versão mais recente | Testes de integração e testes de aceitação |
| US-04 | Detectar JDK; baixar quando ausente; disponibilizar para uso interno | Testes de integração e validação manual por plataforma |
| US-05 | Gerar binários para Windows, Linux e macOS; publicar via GitHub Releases; incluir checksums | Pipeline CI/CD e validação de release |

## Observações
- Esta matriz será atualizada ao longo do desenvolvimento do projeto.
- O objetivo deste documento é facilitar a verificação de cobertura entre requisitos, backlog, documentação, implementação e testes.
- A situação de cada requisito poderá evoluir para: `Planejado`, `Em andamento`, `Implementado` ou `Validado`.
