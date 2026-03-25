## 1.Arquitetura
O sistema foi projetado com base em uma arquitetura em camadas combinada com o padrão Microkernel, garantindo modularidade, extensibilidade e baixo acoplamento.
A solução atende aos principais atributos de qualidade definidos pela ISO/IEC 25010, com destaque para manutenibilidade, portabilidade e confiabilidade.
Além disso, os princípios SOLID foram aplicados para garantir um design robusto, escalável e de fácil evolução.

![](diagramas/imagens/C4_Container.svg)

## 2.Fluxo de execução
![](diagramas/imagens/sequencia.svg)
![](diagramas/imagens/sequenciahttp.svg)

## 3.Linguagens, tecnologias e ferramentas
| Item do Backlog                             | Linguagem | Tecnologias / Ferramentas                        |
| ------------------------------------------- | --------- | ------------------------------------------------ |
| Épico 1 — Assinador                         | Java      | JDK, Maven/Gradle, JUnit                         |
| Épico 2 — Servidor HTTP do Assinador        | Java      | HTTP Server (Java nativo)         |
| Épico 3 — CLI Assinatura                    | Go        | Cobra, net/http                                  |
| Épico 4 — Simulador HubSaúde                | Go        | os/exec, net/http                                |
| Épico 5 — Provisionamento automático do JDK | Go        | net/http, archive/zip, os                        |
| Épico 6 — Qualidade, Testes e Documentação  | Multi     | JUnit, Go Test, Cucumber, GitHub Actions, Cosign |
