# ADR 003 — Modo servidor HTTP como padrão; modo local explícito

## Status
Aceito

## Contexto
As operações `sign` e `verify` podem ser executadas de duas formas:
1. **Local**: o CLI invoca `java -jar assinador.jar` por subprocesso a cada
   operação (sofre cold start da JVM).
2. **Servidor HTTP**: o assinador.jar roda como processo de longa duração e o
   CLI envia requisições HTTP (menor latência por operação).

O critério E2 da especificação exige que o modo servidor seja o padrão e que
o modo local seja ativado explicitamente.

## Decisão
- **Modo servidor é o padrão**: `sign`/`verify` garantem uma instância viva
  (subindo automaticamente se necessário) e operam via HTTP.
- **Modo local é opt-in** via flag `--local`, que invoca o JAR por subprocesso
  e propaga o exit code, separando stdout (resultado) de stderr (diagnóstico).

## Consequências
- **Positivas**: menor latência no caminho comum (sem cold start repetido),
  alinhamento com o critério E2.
- **Negativas**: o servidor pode ficar ocioso consumindo recursos — mitigado
  pelo auto-shutdown por inatividade (`--timeout`, ver US-01.9). Há dois
  caminhos de código (HTTP e subprocess) a manter e testar.

## Referências
- Especificação upstream (tag fixa):
  https://github.com/kyriosdata/runner/tree/v0.1.1