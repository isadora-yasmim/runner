# ADR 002 — Porta padrão 8080 e estratégia de descoberta de instância

## Status
Aceito

## Contexto
No modo servidor, o assinador.jar escuta HTTP numa porta. O CLI precisa
decidir qual porta usar por padrão e como detectar se já existe uma
instância viva — evitando subir processos duplicados e falhando de forma
clara quando a porta está ocupada por outro programa.

## Decisão
- **Porta padrão: 8080**, configurável via `--port`. Escolhida por ser
  convenção amplamente reconhecida para serviços HTTP de desenvolvimento.
- **Descoberta de instância por health check real**, não apenas por "porta
  ocupada": o CLI faz `GET /health` e considera o servidor vivo somente se
  receber `200 OK`. Uma porta ocupada que não responde ao health check é
  tratada como conflito com processo externo, não como instância do
  assinador.

## Consequências
- **Positivas**: idempotência de start (reutiliza instância viva), mensagens
  de erro distinguindo "assinador já rodando" de "porta tomada por outro
  processo".
- **Negativas**: o health check adiciona uma chamada HTTP de latência ao
  start; uma instância em processo de subida pode ainda não responder 200,
  exigindo retry/espera (ver `waitForServer`).

## Referências
- Especificação upstream (tag fixa):
  https://github.com/kyriosdata/runner/tree/v0.1.1