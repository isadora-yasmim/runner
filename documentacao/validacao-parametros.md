# Planejamento — Validação Inicial de Parâmetros

## Contexto

Atualmente, o sistema já possui validações básicas implementadas diretamente nos comandos `SignCommand` e `VerifyCommand`. Essas validações verificam:

- Se os parâmetros obrigatórios foram informados
- Se os valores não estão vazios
- Algumas regras simples (ex: tamanho mínimo do PIN)

Apesar de funcionais, essas validações estão acopladas aos comandos, o que pode dificultar a manutenção e evolução do código.

---

## Objetivo

Refatorar e melhorar a validação inicial de parâmetros, tornando o código mais organizado, reutilizável e alinhado com boas práticas de desenvolvimento.

---

## Estratégia

A abordagem definida é:

- Criar uma classe utilitária dedicada à validação de parâmetros
- Centralizar toda a lógica de validação nessa classe
- Reduzir duplicação de código
- Melhorar clareza e organização dos comandos

---

## Proposta de implementação

### 1. Criar classe `ParameterValidator`

Responsável por:

- Validar parâmetros do comando `sign`
- Validar parâmetros do comando `verify`

Métodos previstos:

- `validateSign(String document, String tokenPin)`
- `validateVerify(String document, String signature)`

---

### 2. Atualizar `SignCommand`

- Remover validações internas
- Utilizar `ParameterValidator`
- Tratar exceções e retornar mensagens de erro ao usuário

---

### 3. Atualizar `VerifyCommand`

- Remover validações internas
- Utilizar `ParameterValidator`
- Manter lógica de simulação separada da validação

---

## Regras de validação (iniciais)

### Sign

- Documento não pode ser nulo ou vazio
- PIN não pode ser nulo ou vazio
- PIN deve ter no mínimo 4 caracteres

### Verify

- Documento não pode ser nulo ou vazio
- Assinatura não pode ser nula ou vazia

---

## Próximos passos

- Implementar a classe `ParameterValidator`
- Refatorar os comandos existentes
- Testar os cenários já existentes
- Garantir que mensagens de erro sejam claras e padronizadas

---

## Observação

Esta etapa corresponde à **validação inicial de parâmetros**, conforme definido no cronograma do projeto (Semana 2).