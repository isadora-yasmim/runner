# Guia de Execução — Servidor HTTP do assinador

## Pré-requisitos

- Java instalado (JDK 11 ou superior)
- Maven instalado

---

# 1. Compilar o projeto

No terminal, dentro da pasta do `assinador`:

```bash
mvn clean compile
```

---

# 2. Iniciar o servidor HTTP

## Porta padrão (8080)

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server"
```

---

## Porta customizada

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server --port 8081"
```

---

# 3. Endpoints disponíveis

## GET /health

Verifica se o servidor HTTP está ativo.

### Exemplo de requisição

```bash
curl http://localhost:8080/health
```

ou no PowerShell:

```bash
curl.exe http://localhost:8080/health
```

---

### Resposta esperada

```json
{
  "success": true,
  "message": "Assinador HTTP ativo",
  "data": {
    "status": "UP"
  }
}
```

---

## POST /sign

Cria uma assinatura simulada.

### Exemplo de requisição

### Linux/macOS

```bash
curl -X POST http://localhost:8080/sign \
-H "Content-Type: application/json" \
-d '{"document":"documento.txt","tokenPin":"1234"}'
```

---

### Windows PowerShell

```bash
curl.exe --% -X POST http://localhost:8080/sign -H "Content-Type: application/json" -d "{\"document\":\"documento.txt\",\"tokenPin\":\"1234\"}"
```

---

### Corpo JSON

```json
{
  "document": "documento.txt",
  "tokenPin": "1234"
}
```

---

### Resposta esperada

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

---

## POST /validate

Valida uma assinatura simulada.

### Exemplo de requisição

### Linux/macOS

```bash
curl -X POST http://localhost:8080/validate \
-H "Content-Type: application/json" \
-d '{"document":"documento.txt","signature":"mock_hash_abc123_base64_encoded_signature_simulated"}'
```

---

### Windows PowerShell

```bash
curl.exe --% -X POST http://localhost:8080/validate -H "Content-Type: application/json" -d "{\"document\":\"documento.txt\",\"signature\":\"mock_hash_abc123_base64_encoded_signature_simulated\"}"
```

---

### Corpo JSON

```json
{
  "document": "documento.txt",
  "signature": "mock_hash_abc123_base64_encoded_signature_simulated"
}
```

---

### Resposta esperada — assinatura válida

```json
{
  "success": true,
  "message": "Assinatura valida.",
  "data": true
}
```

---

### Resposta esperada — assinatura inválida

```json
{
  "success": false,
  "message": "Assinatura invalida ou corrompida.",
  "data": false
}
```

---

# 4. Testes automatizados

Foram implementados testes automatizados para validação dos endpoints HTTP:

- `GET /health`
- `POST /sign`
- `POST /validate`
- tratamento de método inválido (`HTTP 405`)

---

## Executar testes

```bash
mvn test
```

---

## Resultado esperado

```bash
BUILD SUCCESS
```

---

# 5. Encerrar o servidor

No terminal:

```bash
CTRL + C
```

---

# 6. Problemas comuns

## ❌ Porta já em uso

Erro:

```bash
Address already in use
```

### Solução

Usar outra porta:

```bash
--port 8081
```

OU finalizar o processo que está utilizando a porta:

```bash
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## ❌ PowerShell interpretando JSON incorretamente

No Windows PowerShell, recomenda-se utilizar:

```bash
curl.exe --%
```

para evitar problemas de escape de caracteres em requisições JSON.

---

## ❌ Comando `java` não reconhecido

Verifique se o Java está configurado corretamente no `PATH`.

Exemplo:

```text
C:\Program Files\Java\jdk-23\bin
```

---

# 7. Observações

O servidor HTTP permite que o CLI se comunique com o `assinador.jar` via requisições HTTP, evitando o custo de inicialização da JVM a cada execução (*cold start*).

Essa abordagem melhora a reutilização da instância Java, reduz a latência das operações e prepara a aplicação para comunicação entre processos de forma multiplataforma (Windows, Linux e macOS).

A implementação está alinhada com o modo HTTP definido na especificação oficial do Sistema Runner.