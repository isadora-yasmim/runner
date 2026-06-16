package br.go.ses.assinador.commands;

import br.go.ses.assinador.crypto.SignatureToken;
import br.go.ses.assinador.crypto.TokenFactory;
import br.go.ses.assinador.http.AssinadorHttpServer;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.IOException;
import java.util.concurrent.Callable;

@Command(
    name = "server",
    description = "Inicia o assinador.jar em modo servidor HTTP"
)
public class ServerCommand implements Callable<Integer> {

    @Option(
        names = {"-p", "--port"},
        description = "Porta do servidor HTTP. Padrao: 8080"
    )
    private int port = 8080;

    @Option(
        names = {"--real"},
        description = "Usa dispositivo PKCS#11 real em vez do token simulado (padrao: simulado)"
    )
    private boolean real;

    @Option(
        names = {"-T", "--timeout"},
        description = "Minutos de inatividade antes do auto-shutdown. 0 = desativado (padrao)."
    )
    private long timeout = 0L;

    @Override
    public Integer call() {
        try {
            SignatureToken token = TokenFactory.create(real);
            AssinadorHttpServer server = new AssinadorHttpServer(port, token, timeout);
            server.start();

            System.out.println("Servidor HTTP do assinador iniciado com sucesso.");
            System.out.println("Porta: " + port);
            System.out.println("Token: " + token.describe());
            System.out.println("Health check: http://localhost:" + port + "/health");
            if (timeout > 0) {
                System.out.println("Auto-shutdown: " + timeout + " min de inatividade.");
            }
            System.out.println("Pressione CTRL+C para encerrar o servidor.");

            Thread.currentThread().join();

            return 0;

        } catch (IOException e) {
            if (e.getMessage() != null && e.getMessage().contains("Address already in use")) {
                System.err.println("Erro: a porta " + port + " ja esta em uso.");
                System.err.println("Use outra porta com: server --port 8081");
            } else {
                System.err.println("Erro ao iniciar servidor HTTP: " + e.getMessage());
            }
            return 2;

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.err.println("Servidor HTTP interrompido.");
            return 1;
        }
    }
}