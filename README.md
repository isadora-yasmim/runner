# Sistema Runner — CLI de Assinatura Digital Simulada

> Implementação do trabalho prático baseada na especificação disponível em  
> [`kyriosdata/runner @ ab4d353`](https://github.com/kyriosdata/runner/tree/ab4d353)  


---

## O que é este projeto

O **Sistema Runner** é uma ferramenta de linha de comando (CLI) que permite criar e validar
assinaturas digitais simuladas por meio do `assinador.jar`, um componente Java que encapsula
toda a lógica de assinatura e validação.

O CLI (`assinatura`) se comunica com o `assinador.jar` de dois modos:

- **Modo servidor (padrão):** o CLI inicia o `assinador.jar` como servidor HTTP e envia as
  requisições via HTTP. Oferece menor latência por eliminar o cold start do processo Java.
- **Modo local (`--local`):** o CLI invoca o `assinador.jar` diretamente via `java -jar` a
  cada comando. Útil quando não há um servidor em execução.

---

## Pré-requisitos

| Ferramenta | Versão mínima | Verificação |
|---|---|---|
| Go | 1.21 | `go version` |
| JDK | 21 | `java -version` |
| Maven | 3.9 | `mvn -version` |

> **Nota:** o JDK deve estar no `PATH`. O CLI verifica a presença do Java antes de
> qualquer operação e emite mensagem orientativa caso esteja ausente.

---

## Como compilar

### 1. Compilar o assinador.jar (Java)

```bash
cd projetos/assinador-java
mvn clean package
```

O artefato gerado estará em `target/assinador.jar`.

### 2. Compilar o CLI (Go)

```bash
cd projetos/assinatura-cli
go build -o assinatura ./...
```

Para compilar com versão e SHA do commit embutidos (recomendado para releases):

```bash
go build \
  -ldflags "-X cmd.Version=$(git describe --tags) -X cmd.Commit=$(git rev-parse --short HEAD)" \
  -o assinatura \
  ./...
```

### 3. Compilar tudo de uma vez

```bash
make build
```

> O `Makefile` executa o build do JAR e do CLI em sequência e verifica os pré-requisitos.

---

## Como executar

### Exibir a versão do CLI

```bash
assinatura version
```

### Criar uma assinatura digital (modo servidor — padrão)

```bash
assinatura sign --document contrato.pdf --token-pin 1234
```

O CLI inicia o `assinador.jar` automaticamente na porta padrão (`8080`) caso não haja
instância em execução, e reutiliza a instância existente nas próximas chamadas.

### Criar uma assinatura digital (modo local)

```bash
assinatura sign --document contrato.pdf --token-pin 1234 --local
```

### Validar uma assinatura

```bash
assinatura verify \
  --document contrato.pdf \
  --signature mock_hash_abc123_base64_encoded_signature_simulated
```

### Gerenciar o servidor HTTP do assinador

```bash
# Iniciar o servidor na porta padrão (8080)
assinatura start

# Iniciar em porta alternativa
assinatura start --port 9090

# Verificar se o servidor está ativo
assinatura status

# Encerrar o servidor
assinatura stop

# Encerrar servidor em porta alternativa
assinatura stop --port 9090
```

### Porta padrão

A porta padrão é `8080`. Todos os comandos aceitam `--port` para sobrescrever.

```bash
assinatura sign --document doc.pdf --token-pin 1234 --port 9090
```

---

## Saída dos comandos

Todos os resultados são emitidos em JSON estruturado para o `stdout`. Diagnósticos
e erros vão para o `stderr`. Exemplo de resposta de sucesso:

```json
{
  "success": true,
  "message": "Assinatura criada com sucesso (Simulacao)",
  "data": {
    "signatureHash": "mock_hash_abc123_base64_encoded_signature_simulated",
    "algorithm": "SHA256withRSA"
  }
}
```

### Códigos de saída

| Código | Significado |
|---|---|
| `0` | Operação concluída com sucesso |
| `1` | Erro causado pelo usuário (parâmetros inválidos, assinatura inválida) |
| `2` | Erro de sistema (JVM ausente, JAR não encontrado, servidor inacessível) |

---

## Como executar os testes

### Testes unitários do assinador.jar

```bash
cd projetos/assinador-java
mvn test
```

### Relatório de cobertura (Java)

```bash
mvn verify
# Relatório disponível em: target/site/jacoco/index.html
```

### Testes de integração e de contrato CLI ↔ JAR (Go)

```bash
cd projetos/assinatura-cli
go test ./... -v
```

### Relatório de cobertura (Go)

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Executar todos os testes de uma vez

```bash
make test
```

---

## Estrutura do repositório

```
.
├── projetos/
│   ├── assinador-java/          # assinador.jar — componente Java
│   │   ├── src/
│   │   │   ├── main/java/br/go/ses/assinador/
│   │   │   │   ├── commands/    # SignCommand, VerifyCommand, ServerCommand
│   │   │   │   ├── http/        # AssinadorHttpServer
│   │   │   │   ├── model/       # ResponseOutput, SignatureData
│   │   │   │   └── util/        # ParameterValidator, JsonUtil
│   │   │   └── test/java/       # Testes unitários e de integração
│   │   └── pom.xml
│   └── assinatura-cli/          # CLI Go
│       ├── cmd/
│       │   ├── root.go          # Comando raiz e flags globais
│       │   ├── sign.go          # Subcomando sign
│       │   ├── verify.go        # Subcomando verify
│       │   ├── start.go         # Subcomando start
│       │   ├── stop.go          # Subcomando stop
│       │   ├── status.go        # Subcomando status
│       │   ├── version.go       # Subcomando version
│       │   └── server_manager.go # Lógica de ciclo de vida do servidor
│       └── main.go
├── docs/
│   └── adr/                     # Architecture Decision Records
│       ├── 001-cobra-cli.md
│       ├── 002-porta-padrao.md
│       └── 003-modo-servidor-default.md
├── .github/
│   └── workflows/
│       ├── ci.yml               # Build + testes em Ubuntu e Windows
│       └── release.yml          # Binários multiplataforma + Cosign
├── .gitignore
├── .gitattributes
├── LICENSE
├── Makefile
└── README.md
```

---

## Variáveis de ambiente

| Variável | Descrição | Padrão |
|---|---|---|
| `ASSINADOR_JAR` | Caminho absoluto para o `assinador.jar` | Detectado automaticamente |
| `ASSINADOR_PORT` | Porta padrão do servidor HTTP | `8080` |

---

## Arquitetura

O projeto segue um modelo de dois processos: o CLI (Go) orquestra o `assinador.jar` (Java),
que é a **autoridade única** de validação de parâmetros e execução das operações de assinatura.

```
Usuário → CLI (Go) ──HTTP──▶ assinador.jar (HTTP server)
                   ──subprocess──▶ assinador.jar (modo local)
```

Consulte os diagramas em [`docs/`](./docs/) para a visão C4 e os fluxos de sequência.

---

## Como contribuir

1. Abra uma _issue_ descrevendo o problema ou a melhoria, referenciando a história de
   usuário correspondente (`US-XX.Y`).
2. Crie um _branch_ com nome `feat/US-01-2-parsing-cli` ou `fix/jar-path-absoluto`.
3. Faça commits atômicos seguindo [Conventional Commits](https://www.conventionalcommits.org/):
   `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`.
4. Abra um Pull Request ligado à _issue_. O CI deve estar verde (build + testes em Ubuntu
   e Windows) antes do merge.

---

## Status atual

| Componente | Situação |
|---|---|
| `assinador.jar` — comandos sign/verify (modo local) | ✅ Implementado |
| `assinador.jar` — servidor HTTP (sign, verify, health, stop) | ✅ Implementado |
| CLI Go — sign/verify via HTTP | ✅ Implementado |
| CLI Go — start/stop/status do servidor | ✅ Implementado |
| CLI Go — modo local (`--local`) | 🔧 Em andamento |
| Provisionamento automático do JDK | 🔧 Em andamento |
| Simulador HubSaúde (start/stop/status) | 📋 Planejado |
| CI/CD — GitHub Actions multiplataforma | 📋 Planejado |
| Releases com checksums SHA256 + Cosign | 📋 Planejado |

---

## Licença

Este projeto está licenciado sob a [Apache License 2.0](./LICENSE), compatível com as
dependências utilizadas (picocli, Jackson, Cobra).