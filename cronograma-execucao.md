## 📅 Cronograma de Desenvolvimento

### 🟦 Semana 1 — Planejamento
- Definição da arquitetura
- Escolha de tecnologias
- Backlog inicial

### 🟦 Semana 2 — Assinador (base)
- Implementação da simulação de assinatura e validação
- Validação inicial de parâmetros

### 🟦 Semana 3 — Assinador (refinamento)
- Validação completa de parâmetros
- Tratamento de erros
- Testes unitários iniciais

### 🟦 Semana 4 — Modo HTTP
- Implementação do servidor HTTP no `assinador.jar`
- Criação dos endpoints (`/sign`, `/verify`)

### 🟦 Semana 5 — CLI (integração direta)
- Implementação da CLI
- Integração com `assinador.jar` via linha de comando

### 🟦 Semana 6 — CLI (HTTP)
- Integração da CLI com o modo servidor (HTTP)
- Melhoria da saída e usabilidade

### 🟦 Semana 7 — Simulador
- Download automático do `simulador.jar`
- Comandos de iniciar, parar e status

### 🟦 Semana 8 — Provisionamento do JDK
- Detecção do Java instalado
- Download e configuração automática do JDK

### 🟦 Semana 9 — Testes
- Testes de integração
- Cobertura de cenários de erro

### 🟦 Semana 10 — Documentação
- Manual de uso e guia de instalação

### 🟦 Semana 11 — Build e Release
- Geração de binários multiplataforma
- Versionamento (SemVer) e GitHub Releases

### 🟦 Semana 12 — Segurança e Finalização
- Assinatura dos artefatos (Cosign)
- Revisão geral do projeto
