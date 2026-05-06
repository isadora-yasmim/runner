# Guia de Execução — CLI de Assinatura

## 1. Acessar o diretório do projeto

No terminal, navegue até a pasta do CLI:

```bash
cd runner/simulador/cli
```

---

## 2. Criar uma assinatura simulada

```bash
go run . sign -d documento.txt -t 1234
```

### ✔ Saída esperada:

```text
✔ Assinatura criada com sucesso

→ Hash: mock_hash_abc123_base64_encoded_signature_simulated
→ Algoritmo: SHA256withRSA
```

---

## 3. Validar uma assinatura

```bash
go run . verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated
```

### ✔ Saída esperada:

```text
✔ Assinatura válida
```

---

## 4. Validar assinatura inválida

```bash
go run . verify -d documento.txt -s assinatura_errada
```

### Saída esperada:

```text
❌ Assinatura inválida
→ Assinatura invalida ou corrompida.
```

---

## 5. Caminho do `assinador.jar`

Por padrão, o CLI espera o `.jar` em:

```text
../assinador/target/assinador-1.0-SNAPSHOT.jar
```

Caso esteja em outro local, use:

```bash
go run . sign -d documento.txt -t 1234 --jar /caminho/para/assinador.jar
```

---

## Observação importante

O CLI executa internamente:

```bash
java -jar assinador.jar sign -d documento.txt -t 1234
```

Ou seja:

* ✔ Não é necessário executar Java manualmente
* ✔ O CLI abstrai toda a execução
* ✔ A saída já é formatada de forma amigável

---

## Pré-requisitos

* Go instalado
* Java (JDK 21+) disponível no PATH
* `assinador.jar` previamente compilado

