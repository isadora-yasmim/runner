# Guia de Execução — CLI de Assinatura

## Pré-requisitos

- Go instalado
- Java (JDK 21+) disponível no `PATH`
- Maven instalado
- `assinador.jar` previamente compilado

---

## 1. Acessar o diretório do CLI

No terminal, navegue até a pasta do CLI:

```bash
cd runner/simulador/cli
```

---

## 2. Compilar o `assinador.jar`

Antes de utilizar o CLI com integração HTTP, compile o `assinador.jar`.

Na pasta do assinador:

```bash
cd ../assinador
mvn clean package
```

O `.jar` será gerado em:

```text
target/assinador-1.0-SNAPSHOT.jar
```

Depois, retorne para a pasta do CLI:

```bash
cd ../cli
```

---

## 3. Verificar versão do CLI

```bash
go run . version
```

Saída esperada:

```text
assinatura version 0.1.0
```

---

## 4. Iniciar o assinador em modo servidor HTTP

O comando `start` inicia o `assinador.jar` em modo servidor HTTP.

### Porta padrão (8080)

```bash
go run . start
```

### Porta personalizada

```bash
go run . --port 8081 start
```

Saída esperada:

```text
✅ Assinador iniciado na porta 8081.
PID: 12345
```

---

## 5. Verificar status do assinador

### Porta padrão

```bash
go run . status
```

### Porta personalizada

```bash
go run . --port 8081 status
```

Saída esperada quando o servidor está ativo:

```text
✅ Assinador HTTP ativo na porta 8081.
→ O CLI reutilizará esta instância nas próximas operações.
```

Saída esperada quando o servidor não está ativo:

```text
❌ Nenhum assinador HTTP ativo na porta 8081.
→ Os comandos sign e verify tentarão iniciar o assinador automaticamente.
```

---

## 6. Criar uma assinatura simulada

### Porta padrão

```bash
go run . sign -d documento.txt -t 1234
```

### Porta personalizada

```bash
go run . --port 8081 sign -d documento.txt -t 1234
```

Saída esperada:

```text
✔ Assinatura criada com sucesso

→ Hash: mock_hash_abc123_base64_encoded_signature_simulated
→ Algoritmo: SHA256withRSA
```

---

## 7. Validar uma assinatura

### Assinatura válida

```bash
go run . verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
```

Com porta personalizada:

```bash
go run . --port 8081 verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
```

Saída esperada:

```text
✔ Assinatura válida
```

---

### Assinatura inválida

```bash
go run . verify -d documento.txt -s assinatura_errada
```

Com porta personalizada:

```bash
go run . --port 8081 verify -d documento.txt -s assinatura_errada
```

Saída esperada:

```text
❌ Assinatura inválida
→ Assinatura invalida ou corrompida.
```

Observação: nesse caso, o CLI pode finalizar com código de erro, pois a assinatura foi considerada inválida.

---

## 8. Encerrar o assinador

### Porta padrão

```bash
go run . stop
```

### Porta personalizada

```bash
go run . --port 8081 stop
```

Saída esperada:

```text
✅ Assinador encerrado na porta 8081.
```

---

## 9. Executar o servidor Java manualmente

Caso deseje iniciar o servidor Java manualmente, execute dentro da pasta do `assinador`:

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server --port 8081"
```

---

## 10. Fluxo da integração HTTP

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
Resposta
```

Endpoints utilizados:

```text
GET  /health
POST /sign
POST /validate
POST /stop
```

---

## 11. Reutilização de instância

Antes de iniciar uma nova instância do `assinador.jar`, o CLI verifica automaticamente:

```text
GET /health
```

Se o servidor já estiver ativo, a instância existente é reutilizada.

Se o servidor não estiver ativo, os comandos `sign` e `verify` tentam iniciar o assinador automaticamente.

---

## 12. Problemas comuns

### Comando `java` não reconhecido

Verifique se o Java está disponível no `PATH`.

No Windows, o caminho pode ser semelhante a:

```text
C:\Program Files\Java\jdk-23\bin
```

Teste com:

```bash
java -version
```

---

### Arquivo `.jar` não encontrado

Se o CLI não conseguir iniciar o assinador, verifique se o `.jar` foi gerado:

```bash
cd ../assinador
mvn clean package
```

---

### Porta já em uso

Se a porta padrão estiver ocupada, utilize outra porta:

```bash
go run . --port 8081 start
```