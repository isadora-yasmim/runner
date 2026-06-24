# Sistema Runner — CLI de Assinatura Digital Simulada

> Implementação do trabalho prático baseada na especificação disponível em
> [`kyriosdata/runner @ ab4d353`](https://github.com/kyriosdata/runner/tree/ab4d353)

---

## O que é este projeto

O **Sistema Runner** é uma ferramenta de linha de comando (CLI) que permite criar e validar
assinaturas digitais simuladas por meio do `assinador.jar`, um componente Java que encapsula
toda a lógica de assinatura e validação.

O projeto é composto principalmente por dois módulos:

* **Assinador Java (`simulador/assinador`)**: gera o `assinador.jar`, executa comandos de assinatura/validação e disponibiliza endpoints HTTP.
* **CLI Go (`simulador/cli`)**: fornece comandos de linha de comando para interagir com o assinador, iniciar/parar o servidor HTTP e executar operações de assinatura.

O CLI (`assinatura`) se comunica com o `assinador.jar` de dois modos:

- **Modo servidor (padrão):** o CLI inicia o `assinador.jar` como servidor HTTP e envia as
  requisições via HTTP. Oferece menor latência por eliminar o cold start do processo Java.
- **Modo local (`--local`):** o CLI invoca o `assinador.jar` diretamente via `java -jar` a
  cada comando. Útil quando não há um servidor em execução.

---

## Pré-requisitos

| Ferramenta |       Versão mínima recomendada | Verificação     |
| ---------- | ------------------------------: | --------------- |
| Go         | 1.26.2 | `go version`    |
| JDK        |                              21 | `java -version` |
| Maven      |                             3.9 | `mvn -version`  |
| Git        |            versão atual estável | `git --version` |

> **Provisionamento automático do JDK:** caso o JDK não esteja instalado ou
> seja inferior à versão 21, o CLI detecta automaticamente e baixa o JDK 21
> via [Adoptium](https://adoptium.net), armazenando em `~/.hubsaude/jdk/`.
> O download ocorre apenas uma vez — nas execuções seguintes o JDK provisionado
> é reutilizado. Caso o download falhe, o CLI exibe orientação para instalação manual.
---

## Como obter o projeto

Clone o repositório e acesse a raiz do projeto:

```bash
git clone https://github.com/isadora-yasmim/runner.git
cd runner
```

Todos os comandos abaixo assumem que você está na raiz do repositório.

---

## Como compilar

### 1. Compilar o `assinador.jar` Java

A partir da raiz do repositório:

```bash
cd simulador/assinador
mvn clean package
```

O artefato será gerado em `simulador/assinador/target/assinador-1.0-SNAPSHOT.jar`.

Depois, retorne para a raiz do projeto:

```bash
cd ../..
```

---

### 2. Compilar o CLI Go

A partir da raiz do repositório:

```bash
cd simulador/cli
go build
```

No Windows, esse comando gera um executável como `simulador.exe`

Também é possível gerar um binário com nome explícito:

```bash
go build -o assinatura
```

No Windows, o executável gerado será: `assinatura.exe`

Depois, retorne para a raiz:

```bash
cd ../..
```

---

### 3. Compilar com versão e SHA do commit (opcional)

O workflow de release utiliza `ldflags` para injetar informações de versão e rastreabilidade nos binários gerados.

Esse processo é realizado automaticamente pelo pipeline de release e não é necessário para uso ou desenvolvimento local.

---

### 4. Build único

O uso de um `Makefile` para executar o build completo do projeto ainda está previsto como melhoria de reprodutibilidade.

Enquanto o `Makefile` não estiver disponível, use os comandos manuais que foram descritos acima.

---

## Como executar

### Usando `go run` durante o desenvolvimento

A forma mais simples para desenvolvimento é executar o CLI diretamente com `go run`.

Primeiro, gere o `assinador.jar`:

```bash
cd simulador/assinador
mvn clean package
```

Depois, acesse o CLI:

```bash
cd ../cli
```

Agora execute os comandos:

```bash
go run . version
```

Exemplo de saída: `assinatura version 0.1.0`

---

```bash
go run . status
```

Exemplo de saída quando o servidor não está em execução: `❌ Nenhum assinador HTTP ativo na porta 8080. → Os comandos sign e verify tentarão iniciar o assinador automaticamente.`

Exemplo de saída quando o servidor já está ativo: `✅ Assinador HTTP ativo na porta 8080. → O CLI reutilizará esta instância nas próximas operações.`

---

```bash
go run . start
```

Exemplo de saída: `✅ Assinador iniciado na porta 8080. PID: 1252`

> O valor do PID varia a cada execução.

---

```bash
go run . sign -d documento.txt -t 1234
```

Exemplo de saída: `✔ Assinatura criada com sucesso → Hash: mock_hash_abc123_base64_encoded_signature_simulated → Algoritmo: SHA256withRSA `

---

```bash
go run . verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
```

Exemplo de saída para uma assinatura válida: `✔ Assinatura válida`

Exemplo de saída para uma assinatura inválida: `❌ Assinatura inválida → Assinatura invalida ou corrompida.`

---

```bash
go run . stop
```

Exemplo de saída: `✅ Assinador encerrado na porta 8080.`

---

### Usando o binário compilado

Dentro de `simulador/cli`, compile o binário.

#### Linux/macOS

```bash
go build -o assinatura
```

Executando:

```bash
./assinatura version
./assinatura status
./assinatura start
./assinatura sign -d documento.txt -t 1234
./assinatura verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
./assinatura stop
```

---

#### Windows (PowerShell)

Compile o executável:

`go build -o assinatura.exe`

Executando:

```powershell
.\assinatura.exe version
.\assinatura.exe status
.\assinatura.exe start
.\assinatura.exe sign -d documento.txt -t 1234
.\assinatura.exe verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
.\assinatura.exe stop
```

Exemplo de saída: `assinatura version 0.1.0`

---

## Endpoints HTTP do assinador

O CLI se comunica com o `assinador.jar` por HTTP.

Endpoints atualmente utilizados:

```text
GET  /health
POST /sign
POST /validate
POST /stop
```

O comando `verify` do CLI utiliza internamente o endpoint:

```text
POST /validate
```

---

## Como executar os testes

### Testes do módulo Java

Todos os comandos desta seção assumem que você está na raiz do repositório (assinador/).

A partir da raiz:

```bash
cd simulador/assinador
mvn test
```

Para executar o ciclo completo de verificação Maven, incluindo relatório JaCoCo:

```bash
mvn verify
```

O relatório de cobertura Java é gerado em:

```text
simulador/assinador/target/site/jacoco/index.html
```

Depois, retorne para a raiz:

```bash
cd ../..
```

---

### Testes do módulo Go

A partir da pasta `simulador/assinador`, acesse o módulo CLI:

```bash
cd ../cli
```

Execute os testes Go:

```bash
go test ./... -v
```

Atualmente, o módulo Go ainda pode exibir:

```text
[no test files]
```

em alguns pacotes, indicando que os testes do CLI ainda estão em evolução.

A geração de relatório de cobertura Go será usada quando houver testes implementados no módulo CLI. Nesse caso, o comando esperado será:

```bash
go test ./... -coverprofile=coverage.out
```

E, quando o arquivo `coverage.out` for gerado, será possível visualizar a cobertura com:

```bash
go tool cover -html=coverage.out
```


---

## Integração contínua e release

O projeto possui workflows em:

```text
.github/workflows/
├── ci.yml
└── release.yml
```

### CI

O workflow de CI executa build e testes em Ubuntu e Windows para os módulos Java e Go.

Também valida mensagens de commit seguindo o padrão Conventional Commits em Pull Requests.

### Release

O workflow de release é acionado por tags no formato:

```text
vX.Y.Z
```

Ele está configurado para gerar binários multiplataforma, checksums SHA256 e assinaturas com Cosign.

A infraestrutura de release está implementada, mas a validação completa de uma release real depende da criação e publicação de uma tag.

---

### Verificação de integridade e autenticidade

Cada release publica:

- Binários multiplataforma (`linux`, `windows` e `darwin`)
- Arquivo de checksums SHA256
- Assinaturas Cosign (`.sig`)
- Certificados Cosign (`.pem`)

#### Verificar integridade (SHA256)

Após baixar os arquivos da release:

```bash
sha256sum --check checksums-vX.Y.Z.txt
```

O comando deve indicar `OK` para todos os artefatos.

#### Verificar autenticidade (Cosign)

Exemplo para o binário Linux:

```bash
cosign verify-blob \
  --certificate assinatura-vX.Y.Z-linux-amd64.pem \
  --signature assinatura-vX.Y.Z-linux-amd64.sig \
  --certificate-identity-regexp "https://github.com/isadora-yasmim/runner/.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  assinatura-vX.Y.Z-linux-amd64
```

Para Windows ou macOS, substitua os nomes dos arquivos pelo artefato correspondente da release.

A verificação confirma que o artefato foi gerado pelo workflow oficial do projeto utilizando assinatura keyless via GitHub Actions OIDC e registrado no Transparency Log da Sigstore.

---

## Estrutura do repositório

```
runner
├── .github/
│   └── workflows/
│       ├── ci.yml                         # CI: build e testes em Ubuntu e Windows
│       └── release.yml                    # Release: binários multiplataforma, SHA256 e 
├── .vscode/                               # Configurações locais/workspace do VS Code
├── diagramas/
│   ├── imagens/                           # SVGs gerados a partir dos diagramas
│   ├── c4.puml
│   ├── sequencia.puml
│   └── sequenciahttp.puml
├── documentacao/
│   ├── assinador/                         # Documentação do componente Java
│   ├── cli/                               # Documentação do CLI Go
│   ├── servidor/                          # Documentação do modo servidor HTTP
│   ├── servidor http/                     # Pasta legada, em processo de organização
│   ├── definicoes.md
│   ├── design.md
│   ├── especificacao.md
│   └── pipeline.md                        # Documentação dos workflows de CI/CD e release
├── gerenciamento/
│   ├── backlog.md
│   ├── cronograma-execucao.md
│   ├── matriz-rastreabilidade.md
│   └── plano-revisitado-v2.md
├── simulador/
│   ├── assinador/                         # Módulo Java responsável pelo assinador.jar
│   │   ├── src/
│   │   │   ├── main/java/br/go/ses/assinador/
│   │   │   │   ├── Main.java
│   │   │   │   ├── commands/
│   │   │   │   │   ├── ServerCommand.java
│   │   │   │   │   ├── SignCommand.java
│   │   │   │   │   └── VerifyCommand.java
│   │   │   │   ├── http/
│   │   │   │   │   └── AssinadorHttpServer.java
│   │   │   │   ├── model/
│   │   │   │   │   ├── ResponseOutput.java
│   │   │   │   │   └── SignatureData.java
│   │   │   │   └── util/
│   │   │   │       ├── JsonUtil.java
│   │   │   │       └── ParameterValidator.java
│   │   │   └── test/java/br/go/ses/assinador/
│   │   │       ├── commands/
│   │   │       ├── http/
│   │   │       ├── model/
│   │   │       └── util/
│   │   ├── pom.xml
│   │   └── target/                        # Gerado localmente pelo Maven e ignorado pelo Git
│   └── cli/                               # Módulo Go responsável pelo CLI assinatura
│       ├── cmd/
│       │   ├── root.go                    # Comando raiz e flags globais
│       │   ├── server_manager.go          # Gerenciamento do ciclo de vida do servidor
│       │   ├── sign.go                    # Subcomando sign
│       │   ├── start.go                   # Subcomando start
│       │   ├── status.go                  # Subcomando status
│       │   ├── stop.go                    # Subcomando stop
│       │   ├── verify.go                  # Subcomando verify
│       │   └── version.go                 # Subcomando version
│       ├── internal/
│       │   └── executor/
│       │       └── java_executor.go       # Execução auxiliar de processos Java
│       ├── go.mod
│       ├── go.sum
│       ├── golangci.yml
│       ├── LICENSE                        # Licença atual do módulo CLI
│       └── main.go
├── .gitattributes
├── .gitignore
├── README.md
├── geraimagens.bat                        # Script auxiliar para geração de imagens dos diagramas
└── geraimagens.sh                         # Script auxiliar para geração de imagens dos diagramas

```


`target/` pode aparecer localmente porque o Maven gerou, mas é ignorado pelo Git. Já `assinatura.exe` e `simulador.exe` são binários locais gerados nos seus testes, então eu **não colocaria na estrutura oficial do README**, porque eles não fazem parte do repositório versionado.



## Variáveis de ambiente

| Variável         | Descrição                      | Status                  |
| ---------------- | ------------------------------ | ----------------------- |
| `ASSINADOR_JAR`  | Caminho para o `assinador.jar` | Planejada / em evolução |
| `ASSINADOR_PORT` | Porta padrão do servidor HTTP  | Planejada / em evolução |

Atualmente, a porta é configurada principalmente pela flag global:

```bash
--port
```

Exemplo:

```bash
go run . --port 8081 start
```

---

## Arquitetura

O projeto segue um modelo com dois componentes principais:

```text
Usuário
   │
   ▼
CLI (Go)
   ├── HTTP ─────────────▶ assinador.jar (modo servidor HTTP)
   │
   └── subprocess/java ─▶ assinador.jar (modo local)
```

No modo atual, o CLI gerencia o ciclo de vida do `assinador.jar` em modo servidor HTTP:

1. verifica se há um servidor ativo via `/health`;
2. inicia o `assinador.jar` automaticamente se necessário;
3. envia requisições para `/sign` ou `/validate`;
4. reutiliza a instância ativa nas próximas operações;
5. encerra o servidor via `/stop` quando solicitado.

O modo local com `--local` está previsto, mas ainda está em desenvolvimento.

---

## Como contribuir

1. Abra uma issue descrevendo o problema, melhoria ou requisito.
2. Relacione a issue à história de usuário ou critério correspondente.
3. Crie uma branch com nome claro, por exemplo:

   ```text
   feat/modo-local
   fix/jar-path
   docs/atualiza-readme
   ```
4. Faça commits atômicos seguindo Conventional Commits:

   ```text
   feat:
   fix:
   test:
   docs:
   refactor:
   chore:
   ci:
   ```
5. Abra um Pull Request ligado à issue.
6. Aguarde o CI ficar verde.
7. Solicite revisão de pelo menos um membro do time antes do merge.

---

## Status atual

| Componente                                         | Situação                                       |
| -------------------------------------------------- | ---------------------------------------------- |
| `assinador.jar` — comandos sign/verify             | ✅ Implementado                                 |
| `assinador.jar` — servidor HTTP                    | ✅ Implementado                                 |
| Endpoints `/health`, `/sign`, `/validate`, `/stop` | ✅ Implementado                                 |
| CLI Go — `sign` e `verify` via HTTP                | ✅ Implementado                                 |
| CLI Go — `start`, `stop` e `status`                | ✅ Implementado                                 |
| Detecção inicial do Java via PATH                  | ✅ Implementado                                 |
| Validação da existência do `assinador.jar`         | ✅ Implementado                                 |
| CI — GitHub Actions multiplataforma                | ✅ Implementado                                 |
| Release workflow com SHA256 e Cosign               | ✅ Implementado                                 |
| CLI Go — modo local `--local`                      | ✅ Implementado                                 |
| Provisionamento automático do JDK                  | ✅ Implementado                                 |
| Simulador HubSaúde                                 | ✅ Implementado                                 |
| Testes Go de contrato CLI ↔ JAR                    | 📋 Planejado                                   |
| Makefile / build único                             | ✅ Implementado                                 |

---

<!-- 🔧 Em andamento   📋 Planejado  ✅ Implementado -->
S
## Licença

O módulo CLI possui licença Apache License 2.0, disponível em:

```text
simulador/cli/LICENSE
```

A licença é compatível com as dependências utilizadas pelo projeto, incluindo Cobra, Picocli e Jackson.
