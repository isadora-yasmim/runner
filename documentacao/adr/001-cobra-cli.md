# ADR 001 — Uso do Cobra como framework de CLI em Go

## Status
Aceito

## Contexto
O Sistema Runner expõe sua funcionalidade por linha de comando (comandos
`sign`, `verify`, `start`, `stop`, `status`, `version`). Precisávamos de uma
forma consistente de declarar comandos, subcomandos, flags, mensagens de
ajuda e versionamento, sem reinventar parsing de argumentos.

Alternativas consideradas:
- **`flag` (stdlib)**: zero dependências, mas exige construir manualmente a
  árvore de subcomandos, help e validação — verboso e propenso a erro.
- **urfave/cli**: capaz, porém com modelo de configuração menos difundido na
  comunidade Go do que o Cobra.
- **Cobra**: padrão de fato em ferramentas Go (kubectl, gh, hugo), com
  geração de help, subcomandos e integração natural com flags persistentes.

## Decisão
Adotar **spf13/cobra** como framework de CLI, inicializado via `cobra-cli`.
As flags persistentes (ex.: `--port`) ficam no comando raiz; flags locais
(ex.: `--timeout`, `--document`) ficam nos respectivos subcomandos.

## Consequências
- **Positivas**: help automático, estrutura de subcomandos clara, padrão
  reconhecido por quem mantém o projeto.
- **Negativas**: adiciona uma dependência externa e o scaffolding inicial do
  `cobra-cli` traz código boilerplate (ex.: flag `toggle`) que precisa ser
  limpo manualmente.
- O `--help` e o `--version` passam a ser contrato observável do CLI e devem
  ser cobertos por testes.

## Referências
- Especificação upstream (tag fixa):
  https://github.com/kyriosdata/runner/tree/v0.1.1