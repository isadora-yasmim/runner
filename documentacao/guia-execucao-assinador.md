# Guia de Execucao - Assinador

Este documento descreve como compilar e executar o modulo `assinador` do Sistema Runner.

## 1. Acessar a pasta do projeto

No terminal, entre na pasta do modulo:


```powershell
cd simulador\assinador
```

## 2. Compilar o projeto

Para compilar o projeto com Maven:

```powershell
mvn compile
```
Se estiver tudo correto, o terminal exibira BUILD SUCCESS.

## 3. Exibir ajuda do programa

Para visualizar os comandos disponiveis:
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=--help"
```
## 4. Executar sem subcomando

Para verificar o comportamento padrao do programa:
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main"
```
Saida esperada:

`Use 'assinador --help' para ver os comandos disponiveis.`

## 5. Comando sign
### 5.1 Assinatura valida
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=sign -d documento.txt -t 1234"
```
Saida esperada:

`success: true`
mensagem de sucesso
dados simulados da assinatura

### 5.2 PIN invalido
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=sign -d documento.txt -t 12"
```
Saida esperada:

`success: false`
mensagem informando que o PIN deve ter pelo menos 4 caracteres

### 5.3 Documento vazio
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=sign -d= -t 1234"
```
Saida esperada:

`success: false`
mensagem informando que o documento nao pode estar vazio

## 6. Comando verify

### 6.1 Assinatura valida

```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=verify -d documento.txt -s mock_hash_abc123_base64_encoded_signature_simulated"
```

Saida esperada:

`success: true`
mensagem informando que a assinatura e valida

### 6.2 Assinatura invalida
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=verify -d documento.txt -s assinatura_errada"
```
Saida esperada:

`success: false`
mensagem informando que a assinatura e invalida ou corrompida

### 6.3 Documento vazio
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=verify -d= -s assinatura_teste"
```
Saida esperada:

`success: false`
mensagem informando que o documento e obrigatorio para a validacao

### 6.4 Assinatura vazia
```powershell
mvn exec:java "-Dexec.mainClass=br.go.ses.assinador.Main" "-Dexec.args=verify -d documento.txt -s="
```
Saida esperada:

`success: false`
mensagem informando que a assinatura e obrigatoria para a validacao

## 7. Observacoes
O modulo assinador implementa apenas simulacao de assinatura digital.
O foco atual esta na validacao inicial de parametros e na organizacao da logica.
A validacao dos parametros foi centralizada na classe ParameterValidator.