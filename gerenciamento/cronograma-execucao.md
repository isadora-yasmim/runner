## 📅 Cronograma de Desenvolvimento (v.2)

### 🟦 Semana 1 — Planejamento
- [x] Definicao da arquitetura
- [x] Escolha de tecnologias
- [x] Backlog inicial

### 🟦 Semana 2 — Assinador (base)
- [x] Simulacao de assinatura e validacao
- [x] Validacao inicial de parametros

### 🟦 Semana 3 — Assinador (refinamento)
- [x] Validacao completa de parametros
- [x] Tratamento de erros
- [x] Testes unitarios iniciais

### 🟦 Semana 4 — CLI (integracao direta) + testes iniciais
- [x] Implementacao da CLI
- [x] Parsing de comandos
- [x] Integracao com `assinador.jar`
- [x] Testes basicos de integracao

### 🟦 Semana 5 — Modo HTTP + melhorias CLI
- [x] Servidor HTTP no `assinador.jar`
- [x] Endpoints `/sign` e `/validate`
- [x] Reaproveitamento da validacao
- [x] Ajustes na CLI

### 🟦 Semana 6 — CLI (HTTP) + gerenciamento de processo
- [x] Integracao CLI com HTTP
- [x] Reutilizacao de instancia
- [x] Inicio/stop do assinador
- [ ] Melhorias de usabilidade

### 🟦 Semana 7 — Simulador + Provisionamento do JDK
- [ ] Download do `simulador.jar`
- [ ] Comandos start/stop/status
- [ ] Deteccao do Java
- [ ] Download automatico do JDK

### 🟦 Semana 8 — Testes e integracao geral
- [ ] Testes de integracao completos
- [ ] Cenarios de erro
- [ ] Ajustes finais no fluxo

### 🟦 Semana 9 — Documentacao + Build
- [ ] Manual de uso
- [ ] Guia de instalacao
- [ ] Documentacao tecnica
- [ ] Build multiplataforma

### 🟦 Semana 10 — Release + Seguranca + Finalizacao
- [ ] GitHub Releases
- [ ] SemVer
- [ ] Checksums
- [ ] Assinatura com Cosign
- [ ] Revisao geral