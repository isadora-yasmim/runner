# Pipeline de CI/CD — Sistema Runner

Este documento descreve os dois workflows de automação do projeto, localizados em `.github/workflows/`. Eles garantem que todo código integrado à branch principal compile, passe nos testes e seja distribuído de forma segura e rastreável.

---

## Visão geral

```
push / pull request         push de tag vX.Y.Z
        │                          │
        ▼                          ▼
  ┌─────────────┐          ┌───────────────┐
  │   ci.yml    │          │ release.yml   │
  │             │          │               │
  │ java        │          │ build-jar     │
  │ go          │          │ build-go      │
  │ lint-go     │          │   ├ linux     │
  │ commitlint  │          │   ├ windows   │
  └─────────────┘          │   └ darwin    │
                           │ publish       │
                           │   ├ SHA256    │
                           │   ├ Cosign    │
                           │   └ Release   │
                           └───────────────┘
```

Os dois workflows são independentes: o CI roda em todo push e PR, enquanto o Release só é acionado quando uma tag semântica é criada.

---

## ci.yml — Integração Contínua

**Arquivo:** `.github/workflows/ci.yml`  
**Acionado por:** push na branch `main` e abertura de Pull Request.  
**Plataformas:** `ubuntu-latest` e `windows-latest` (ambos os módulos rodam nas duas).

Alterações restritas a documentação (`**.md`, `diagramas/**`, `documentacao/**`, `gerenciamento/**`) não disparam o CI, evitando consumo desnecessário de minutos.

### Jobs

#### `java` — Build e testes do módulo Java

Raiz do módulo: `simulador/assinador/`

| Step | O que faz |
|---|---|
| Checkout | Baixa o código |
| Configurar JDK 21 (Temurin) | Instala o Java com cache do Maven |
| `mvn verify` | Compila, executa todos os testes e gera relatório Jacoco |
| Publicar cobertura Jacoco | Salva `target/site/jacoco/` como artifact por 14 dias (só Linux) |
| Publicar resultados Surefire | Salva os XMLs de teste como artifact por 7 dias |

O comando `mvn --batch-mode -e verify` executa as fases `compile → test → package → verify` em sequência. A flag `-e` exibe o stack trace completo em caso de falha, facilitando o diagnóstico.

#### `go` — Build e testes do módulo Go

Raiz do módulo: `simulador/cli/`

| Step | O que faz |
|---|---|
| Checkout | Baixa o código |
| Configurar Go 1.25 | Instala o Go com cache baseado em `go.sum` |
| `go vet` | Análise estática básica — detecta erros comuns |
| `go build` | Compila todos os pacotes |
| `go test` | Roda os testes com detector de race conditions e coleta de cobertura |
| Publicar cobertura | Salva `coverage.html` como artifact por 14 dias (só Linux) |

O flag `-race` ativa o detector de condições de corrida. O flag `-covermode=atomic` garante medição de cobertura correta mesmo com testes concorrentes.

#### `lint-go` — Análise estática avançada

Roda o `golangci-lint` apontando para `simulador/cli/`, que contém o arquivo de configuração `golangci.yml` com as regras do projeto. Separado do job `go` para que o feedback de lint e de testes chegue em paralelo.

#### `commitlint` — Validação de mensagens de commit

Roda **apenas em Pull Requests**. Valida todas as mensagens de commit do PR usando a convenção [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: adicionar endpoint /validate
fix: corrigir timeout no health check
chore: atualizar dependências do Maven
docs: adicionar guia de execução do CLI
test: cobrir cenário de PIN inválido
```

PRs com mensagens fora desse padrão têm o merge bloqueado.

---

## release.yml — Pipeline de Release

**Arquivo:** `.github/workflows/release.yml`  
**Acionado por:** push de tag no formato `vX.Y.Z` (ex.: `v1.0.0`).  
**Produz:** binários multiplataforma, fat JAR, checksums SHA256 e assinaturas Cosign publicados no GitHub Releases.

Para acionar uma release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Tags com sufixo (ex.: `v1.0.0-beta`) são marcadas automaticamente como pré-release.

### Jobs

#### `build-jar` — Compilação do assinador.jar

Raiz do módulo: `simulador/assinador/`

Executa `mvn package` com os testes habilitados — uma release nunca é publicada com testes falhando. O JAR gerado é padronizado para o nome `assinador.jar` e salvo como artifact interno para uso pelo job `publish`.

#### `build-go` — Cross-compilation do CLI

Raiz do módulo: `simulador/cli/`

Compila o CLI para três alvos em paralelo:

| Sistema | Arquitetura | Arquivo gerado |
|---|---|---|
| Linux | amd64 | `assinatura-vX.Y.Z-linux-amd64` |
| Windows | amd64 | `assinatura-vX.Y.Z-windows-amd64.exe` |
| macOS | amd64 | `assinatura-vX.Y.Z-darwin-amd64` |

A versão e o hash curto do commit são injetados no binário via `-ldflags` durante a compilação, de modo que `assinatura version` retorna algo rastreável:

```
assinatura v1.0.0 (abc1234)
```

#### `publish` — Publicação

Executa somente após `build-jar` e `build-go` concluírem com sucesso.

| Step | O que faz |
|---|---|
| Baixar artefatos | Coleta binários e JAR gerados pelos jobs anteriores |
| Gerar checksums SHA256 | Cria `checksums-vX.Y.Z.txt` com hash de cada artefato |
| Instalar Cosign | Prepara a ferramenta de assinatura |
| Assinar com Cosign | Gera `.sig` e `.pem` para cada artefato (keyless via OIDC) |
| Criar GitHub Release | Publica tudo com notas de release geradas automaticamente |

### Verificar integridade de um artefato baixado

**SHA256:**
```bash
sha256sum --check checksums-v1.0.0.txt
```

**Cosign:**
```bash
cosign verify-blob \
  --certificate     assinatura-v1.0.0-linux-amd64.pem \
  --signature       assinatura-v1.0.0-linux-amd64.sig \
  --certificate-identity-regexp "https://github.com/kyriosdata/runner/*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  assinatura-v1.0.0-linux-amd64
```

---

## Arquivos gerados por release

Cada release no GitHub publica os seguintes arquivos:

```
assinatura-vX.Y.Z-linux-amd64
assinatura-vX.Y.Z-linux-amd64.sig
assinatura-vX.Y.Z-linux-amd64.pem
assinatura-vX.Y.Z-windows-amd64.exe
assinatura-vX.Y.Z-windows-amd64.exe.sig
assinatura-vX.Y.Z-windows-amd64.exe.pem
assinatura-vX.Y.Z-darwin-amd64
assinatura-vX.Y.Z-darwin-amd64.sig
assinatura-vX.Y.Z-darwin-amd64.pem
assinador.jar
assinador.jar.sig
assinador.jar.pem
checksums-vX.Y.Z.txt
```

---

## Onde ficam os arquivos no repositório

```
.github/
└── workflows/
    ├── ci.yml        — Integração Contínua
    └── release.yml   — Pipeline de Release

simulador/
├── assinador/        — módulo Java (Maven)
│   └── golangci.yml  — não se aplica (Java não usa golangci)
└── cli/
    └── golangci.yml  — regras do golangci-lint para o módulo Go
```
