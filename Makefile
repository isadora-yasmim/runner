# Makefile — Sistema Runner
#
# Build reproduzível para os dois módulos do projeto:
#   - simulador/assinador  (Java / Maven  -> assinador.jar)
#   - simulador/cli        (Go            -> binário assinatura)
#
# Uso rápido:
#   make            # build completo (jar + binário Go) + testes
#   make build      # apenas compila os dois módulos
#   make test       # apenas roda os testes dos dois módulos
#   make help       # lista todos os alvos
#
# Os mesmos alvos são usados pelo CI (.github/workflows/ci.yml),
# garantindo paridade entre máquina local e pipeline.

# ----------------------------------------------------------------------
# Configuração
# ----------------------------------------------------------------------

# Diretórios dos módulos
JAVA_DIR := simulador/assinador
GO_DIR   := simulador/cli

# Nome fixo do binário Go (coerente com o release multiplataforma)
GO_BIN := assinatura

# Nome fixo do JAR (definido em <finalName>assinador</finalName> no pom.xml).
# O CLI Go procura por este nome — manter sincronizado.
JAR_NAME := assinador.jar
JAR_PATH := $(JAVA_DIR)/target/$(JAR_NAME)

# Versão e commit injetados no binário Go via ldflags (rastreabilidade).
# Fallback para "dev"/"none" quando fora de um repositório git (ex.: tarball).
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -X github.com/isadora-yasmim/simulador/cmd.Version=$(VERSION) \
           -X github.com/isadora-yasmim/simulador/cmd.Commit=$(COMMIT)

# Detecta o comando Maven (mvnw tem prioridade se existir no módulo Java)
MVN := $(if $(wildcard $(JAVA_DIR)/mvnw),./mvnw,mvn)

# Garante .PHONY para todos os alvos (não são arquivos)
.PHONY: all help build build-java build-go test test-java test-go \
        lint vet fmt clean tidy run-server coverage

# ----------------------------------------------------------------------
# Alvo padrão
# ----------------------------------------------------------------------

all: build test ## Build completo + testes (alvo padrão)

help: ## Mostra esta ajuda
	@echo "Sistema Runner — alvos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

# ----------------------------------------------------------------------
# Build
# ----------------------------------------------------------------------

build: build-java build-go ## Compila os dois módulos (jar + binário Go)

build-java: ## Empacota o assinador.jar (fat JAR via maven-shade-plugin)
	cd $(JAVA_DIR) && $(MVN) -B clean package

build-go: ## Compila o binário Go (assinatura) com versão/commit embutidos
	cd $(GO_DIR) && go build -ldflags "$(LDFLAGS)" -o $(GO_BIN) .

# ----------------------------------------------------------------------
# Testes
# ----------------------------------------------------------------------

test: test-java test-go ## Roda os testes dos dois módulos

test-java: ## Testes Java (JUnit) + relatório JaCoCo (target/site/jacoco)
	cd $(JAVA_DIR) && $(MVN) -B verify

test-go: ## Testes Go (precisa do jar; depende de build-java)
test-go: build-java
	cd $(GO_DIR) && go test ./...

coverage: ## Gera relatório de cobertura Go (coverage.html)
coverage: build-java
	cd $(GO_DIR) && go test -coverprofile=coverage.out ./... \
		&& go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório Go: $(GO_DIR)/coverage.html"
	@echo "Relatório Java: $(JAVA_DIR)/target/site/jacoco/index.html"

# ----------------------------------------------------------------------
# Qualidade de código
# ----------------------------------------------------------------------

lint: vet ## Lint do código Go (golangci-lint se disponível; senão, go vet)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(GO_DIR) && golangci-lint run ./...; \
	else \
		echo "golangci-lint não encontrado — execute 'make vet' ou instale a ferramenta"; \
	fi

vet: ## Análise estática nativa do Go
	cd $(GO_DIR) && go vet ./...

fmt: ## Formata o código Go (gofmt)
	cd $(GO_DIR) && gofmt -w .

tidy: ## Sincroniza go.mod/go.sum
	cd $(GO_DIR) && go mod tidy

# ----------------------------------------------------------------------
# Execução / utilidades
# ----------------------------------------------------------------------

run-server: build ## Sobe o assinador em modo servidor HTTP (porta padrão)
	cd $(GO_DIR) && ./$(GO_BIN) start

clean: ## Remove artefatos de build dos dois módulos
	cd $(JAVA_DIR) && $(MVN) -B clean
	cd $(GO_DIR) && rm -f $(GO_BIN) coverage.out coverage.html
