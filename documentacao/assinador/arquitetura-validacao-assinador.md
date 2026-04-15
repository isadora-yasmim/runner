# Arquitetura Atual - validação do Assinador

## Visão Geral

A validação inicial de parâmetros do módulo `assinador` foi organizada de forma a separar responsabilidades entre as classes do sistema.

Atualmente, a estrutura principal dessa parte esta dividida assim:

- `Main`: ponto de entrada da aplicação
- `SignCommand`: executa o fluxo de assinatura simulada
- `VerifyCommand`: executa o fluxo de validação simulada
- `ParameterValidator`: centraliza as regras de validação de entrada
- `JsonUtil`: formata e exibe as respostas em JSON
- `SignatureData`: representa os dados da assinatura simulada
- `ResponseOutput`: representa o formato padrão de resposta

## Responsabilidades

### Main:
A classe `Main` e responsável por iniciar a aplicação e registrar os subcomandos disponíveis no terminal.

### SignCommand:
A classe `SignCommand` recebe os parâmetros do comando `sign`, solicita a validação desses parâmetros ao `ParameterValidator` e, se estiver tudo correto, executa a simulação da criação de assinatura.

### VerifyCommand:
A classe `VerifyCommand` recebe os parâmetros do comando `verify`, solicita a validação ao `ParameterValidator` e, se os parâmetros forem válidos, executa a simulação da verificação da assinatura.

### ParameterValidator:
A classe `ParameterValidator` foi criada para concentrar as regras de validação de parâmetros em um único ponto do sistema.

Ela possui os métodos:
- `validateSign(String document, String tokenPin)`
- `validateVerify(String document, String signature)`

Essa classe melhora a organização do código e evita repetição de validações dentro dos comandos.

### JsonUtil:
A classe `JsonUtil` e usada para padronizar a saída da aplicação em JSON, tanto em casos de sucesso quanto em casos de erro.

## Fluxo Atual

### Fluxo de `sign`
1. O usuário executa o comando `sign`
2. O `SignCommand` recebe os parâmetros
3. O `SignCommand` chama `ParameterValidator.validateSign(...)`
4. Se houver erro, a exceção é capturada e a mensagem e exibida em JSON
5. Se os parâmetros forem válidos, a assinatura simulada e criada
6. O resultado e exibido em JSON

### Fluxo de `verify`
1. O usuário executa o comando `verify`
2. O `VerifyCommand` recebe os parâmetros
3. O `VerifyCommand` chama `ParameterValidator.validateVerify(...)`
4. Se houver erro, a exceção é capturada e a mensagem e exibida em JSON
5. Se os parâmetros forem válidos, a assinatura e comparada com o valor esperado
6. O resultado é exibido em JSON

## Benefícios da Solução Atual

A estrutura atual traz os seguintes benefícios:

- separação de responsabilidades
- menor repetição de código
- validação centralizada
- maior facilidade de manutenção
- melhor legibilidade dos comandos
- base melhor para evolução futura

## Relação com boas práticas

A mudança aproxima o projeto de boas práticas de modularidade e do princípio de responsabilidade única, pois cada classe passa a ter um papel mais claro no sistema.