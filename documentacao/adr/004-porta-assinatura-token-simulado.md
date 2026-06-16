# ADR 004 — Porta `SignatureToken`, recebedor e token simulado

- **Status:** Aceito
- **Data:** 2026-06-03
- **Issue:** #52
- **Histórias:** US-02.1 (assinatura simulada), prepara US-02.5 (PKCS#11)
- **Critérios do professor:** D (acoplamento transporte/domínio/interface), E3 (validação no domínio), E5 (PKCS#11), A (decisão registrada)

## Contexto

A US-02.1 previa uma interface `SignatureService` e uma implementação
`FakeSignatureService`, mas isso nunca foi feito. O material de assinatura
simulado (`mock_hash_abc123…`, `SHA256withRSA`) estava **hardcoded e
duplicado** em três lugares:

- `SignCommand` — `new SignatureData("mock_hash_abc123…", "SHA256withRSA")`;
- `AssinadorHttpServer` — constante `MOCK_SIGNATURE` e a mesma comparação;
- `VerifyCommand` — repetia a comparação contra o hash.

Consequências: o domínio (assinar/validar) estava acoplado ao transporte
(HTTP) e à interface (CLI); não havia ponto único para trocar entre token
real e simulado; a US-02.5 (dispositivo via PKCS#11) não tinha onde encaixar.

## Decisão

Introduzir uma **porta** única que representa o dispositivo de assinatura e
separar claramente os papéis:

- `SignatureToken` (porta): contrato do token/smartcard — `sign`, `verify`,
  `isPresent`, `describe`.
- `SignatureService` (recebedor): depende **apenas** da porta, concentra a
  validação de parâmetros (autoridade única, criterio E3) e delega a operação
  ao token, sem conhecer a implementação concreta.
- `SimulatedToken` (implementação): in-process, determinística, **única fonte
  da verdade** do material simulado. `isPresent()` retorna sempre `true` — do
  ponto de vista do recebedor é indistinguível de um token real, "enganando-o".
- `Pkcs11Token` (implementação): token real via `SunPKCS11`. Entregue agora
  como **esqueleto** (`isPresent()==false`; `sign`/`verify` lançam
  `TokenException` com mensagem orientativa); implementação real e testes com
  SoftHSM2 na Issue #47.
- `TokenFactory`: seleciona a implementação — simulado por padrão, real via
  flag `--real`.
- `TokenException`: erro de **sistema/dispositivo**, distinto de
  `IllegalArgumentException` (erro do **usuário**).

A validação ocorre antes da checagem de presença do dispositivo, garantindo
que erro do usuário nunca seja mascarado por erro de sistema.

## Consequências

**Positivas**
- Domínio desacoplado de transporte e interface; CLI e HTTP delegam ao mesmo
  `SignatureService`.
- Mock deixa de ser duplicado: trocar a simulação exige mudança em um só lugar.
- US-02.5 passa a ser "preencher `Pkcs11Token`" — sem tocar em CLI/HTTP/serviço.
- Erros do usuário e do sistema mapeados de forma consistente:
  - `IllegalArgumentException` → exit code **1** / HTTP **400**;
  - `TokenException` → exit code **2** / HTTP **503**.

**Custos / riscos**
- Mais classes e uma indireção a mais entre interface e operação.
- A mensagem de sucesso ainda diz "(Simulacao)"; quando a #47 habilitar o modo
  real com sucesso, a mensagem deverá ser tornada neutra.

**Compatibilidade**
- O comportamento simulado é preservado (mesmo hash, mesmo algoritmo, mesmos
  status HTTP e exit codes nos caminhos existentes). A suíte Java e Go atual e
  o contrato JSON com o CLI Go permanecem válidos — o CI não quebra.

## Alternativas consideradas

- **Manter o mock inline** com um `if (real)` espalhado por CLI e HTTP:
  rejeitada por espalhar a decisão e manter o acoplamento.
- **Strategy só no CLI**, sem tocar no servidor HTTP: rejeitada porque deixaria
  o mock duplicado no HTTP e dois pontos de manutenção.
