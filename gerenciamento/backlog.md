# Backlog do Projeto Sistema Runner

## Épico 1 — Assinador (assinador.jar)
- [x] Criar projeto Java
- [x] Implementar comando `sign` (simulação de assinatura)
- [x] Implementar comando `verify` (simulação de validação)
- [x] Criar validador de parâmetros
- [x] Padronizar respostas (formato legível e estruturado)
- [x] Criar erros estruturados
- [x] Criar testes unitários

## Épico 2 — Servidor HTTP do Assinador
- [x] Iniciar servidor HTTP no assinador.jar
- [x] Criar endpoint `POST /sign`
- [x] Criar endpoint `POST /verify`
- [x] Criar endpoint `GET /health`
- [x] Testar API com respostas corretas e tratamento de erros

## Épico 3 — CLI Assinatura
- [x] Criar parser de comandos e subcomandos
- [x] Implementar subcomando `sign`
- [x] Implementar subcomando `verify`
- [ ] Implementar escolha de modo (direto ou HTTP)
- [x] Formatar saídas para o usuário
- [x] Testar integração com assinador.jar (modos direto e HTTP)

## Épico 4 — Simulador HubSaúde
- [ ] Consultar release mais recente do simulador no GitHub
- [ ] Baixar jar se versão mais recente não estiver local
- [ ] Verificar versão local e comparar com remota
- [ ] Iniciar simulador
- [ ] Parar simulador
- [ ] Mostrar status do simulador
- [ ] Verificar portas antes de iniciar

## Épico 5 — Provisionamento automático do JDK
- [x] Detectar se Java está instalado
- [ ] Checar versão compatível
- [ ] Baixar JDK por plataforma se necessário
- [ ] Extrair e disponibilizar JDK local
- [ ] Garantir que assinador.jar e simulador.jar usem o Java local

## Épico 6 — Qualidade, Testes e Documentação
- [ ] Criar testes unitários para lógica e validação
- [ ] Criar testes de integração (CLI ↔ assinador.jar ↔ simulador)
- [ ] Criar testes de aceitação baseados nas histórias de usuário
- [ ] Implementar tratamento de erros consistente
- [ ] Criar pipeline CI/CD para build e release multiplataforma
- [ ] Gerar binários Windows, Linux e macOS
- [ ] Criar checksums SHA256
- [ ] Assinar artefatos com Cosign e publicar `.sig` e `.pem`
