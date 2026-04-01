# Backlog do Projeto Sistema Runner

## Épico 1 — Assinador (assinador.jar)
- [ ] Criar projeto Java
- [ ] Implementar comando `sign` (simulação de assinatura)
- [ ] Implementar comando `verify` (simulação de validação)
- [ ] Criar validador de parâmetros
- [ ] Padronizar respostas (formato legível e estruturado)
- [ ] Criar erros estruturados
- [ ] Criar testes unitários

## Épico 2 — Servidor HTTP do Assinador
- [ ] Iniciar servidor HTTP no assinador.jar
- [ ] Criar endpoint `POST /sign`
- [ ] Criar endpoint `POST /verify`
- [ ] Criar endpoint `GET /health`
- [ ] Testar API com respostas corretas e tratamento de erros

## Épico 3 — CLI Assinatura
- [ ] Criar parser de comandos e subcomandos
- [ ] Implementar subcomando `sign`
- [ ] Implementar subcomando `verify`
- [ ] Implementar escolha de modo (direto ou HTTP)
- [ ] Formatar saídas para o usuário
- [ ] Testar integração com assinador.jar (modos direto e HTTP)

## Épico 4 — Simulador HubSaúde
- [ ] Consultar release mais recente do simulador no GitHub
- [ ] Baixar jar se versão mais recente não estiver local
- [ ] Verificar versão local e comparar com remota
- [ ] Iniciar simulador
- [ ] Parar simulador
- [ ] Mostrar status do simulador
- [ ] Verificar portas antes de iniciar

## Épico 5 — Provisionamento automático do JDK
- [ ] Detectar se Java está instalado
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
