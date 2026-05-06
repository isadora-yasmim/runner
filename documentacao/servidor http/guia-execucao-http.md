# Guia de Execução - Servidor HTTP do assinador

## Pré-requisitos

* Java instalado (JDK 11 ou superior)
* Maven instalado

---

## 1. Compilar o projeto

No terminal, dentro da pasta do `assinador`:

```bash
mvn clean compile
```

---

## 2. Iniciar o servidor HTTP

### Porta padrão (8080)

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server"
```

---

### Porta customizada

```bash
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=server --port 8081"
```

---

## 3. Testar o servidor

Acesse no navegador ou via terminal:

```bash
http://localhost:8080/health
```

ou:

```bash
curl http://localhost:8080/health
```

---

## Resposta esperada

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

## 4. Encerrar o servidor

No terminal:

```bash
CTRL + C
```

---

## ⚠️ Problemas comuns

### ❌ Porta já em uso

Erro:

```bash
Address already in use
```

Solução:

* Usar outra porta:

```bash
--port 8081
```

OU

* Finalizar processo:

```bash
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## Observação

O servidor HTTP permite que o CLI se comunique com o `assinador.jar` via requisições HTTP, evitando o custo de inicialização da JVM a cada execução (modo servidor). Isso está alinhado com o modo HTTP definido na especificação do projeto .
