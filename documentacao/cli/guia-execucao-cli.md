# Guia de Execução — CLI de Assinatura

## 1. Acessar o diretório do projeto

No terminal, navegue até a pasta do CLI:

```bash
cd runner/simulador/cli
```

---

## 2. Iniciar o assinador em modo servidor HTTP

Antes de utilizar os comandos `sign` e `verify` via HTTP, inicie o `assinador.jar` em modo servidor.

### Porta padrão (8080)

```bash
go run . start
```

### Porta personalizada

Exemplo usando a porta `8081`:

```bash
go run . start -p 8081
```

### ✔ Saída esperada

```text
✅ Assinador iniciado na porta 8081.
PID: 12345
```

---

## 3. Verificar status do assinador

### Porta padrão

```bash
go run . status
```

### Porta personalizada

```bash
go run . status -p 8081
```

### ✔ Saída esperada

```text
✅ Assinador está em execução na porta 8081.
```

---

## 4. Criar uma assinatura simulada

### Porta padrão

```bash
go run . sign -d documento.txt -t 1234
```

### Porta personalizada

```bash
go run . sign -d documento.txt -t 1234 -p 8081
```

### ✔ Saída esperada

```json
{
  "success": true,
  "message": "Assinatura criada com sucesso (HTTP)",
  "data": {
    "signatureHash": "mock_hash_abc123_base64_encoded_signature_simulated",
    "algorithm": "SHA256withRSA"
  }
}
```

---

## 5. Validar uma assinatura

```bash
go run . verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated -p 8081
```

### ✔ Saída esperada

```json
{
  "success": true,
  "message": "Assinatura valida.",
  "data": true
}
```

---

## 6. Validar assinatura inválida

```bash
go run . verify -d documento.txt -s assinatura_errada -p 8081
```

### ✔ Saída esperada

```json
{
  "success": false,
  "message": "Assinatura invalida ou corrompida.",
  "data": null
}
```

---

## 7. Encerrar o assinador

### Porta padrão

```bash
go run . stop
```

### Porta personalizada

```bash
go run . stop -p 8081
```

### ✔ Saída esperada

```text
✅ Assinador encerrado na porta 8081.
```

---

## 8. Executar o assinador.jar manualmente (modo servidor)

Caso deseje iniciar o servidor Java manualmente:

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server --port 8081"
```

---

## Fluxo da integração HTTP

O CLI se comunica com o `assinador.jar` via HTTP:

```text
CLI (Go)
   ↓
HTTP Request
   ↓
assinador.jar
   ↓
Validação + Simulação
   ↓
Resposta JSON
```

Endpoints disponíveis:

```text
GET  /health
POST /sign
POST /verify
POST /stop
```

---

## Reutilização de instância

Antes de iniciar uma nova instância do `assinador.jar`, o CLI verifica automaticamente:

```text
GET /health
```

Se o servidor já estiver ativo, a instância existente é reutilizada.

---

## Pré-requisitos

* Go instalado
* Java (JDK 21+) disponível no PATH
* Maven instalado
* `assinador.jar` previamente compilado

---

## Compilar o assinador.jar

Na pasta do assinador:

```bash
mvn clean package
```

O `.jar` será gerado em:

```text
../assinador/target/assinador-1.0-SNAPSHOT.jar
```
