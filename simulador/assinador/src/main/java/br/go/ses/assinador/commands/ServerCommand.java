package br.go.ses.assinador.commands;

import br.go.ses.assinador.http.AssinadorHttpServer;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.IOException;
import java.util.concurrent.Callable;

@Command(name = "server", description = "Inicia o assinador.jar em modo servidor HTTP")
public class ServerCommand implements Callable<Integer> {

    @Option(
        names = {"-p", "--port"},
        description = "Porta do servidor HTTP. Padrao: 8080"
    )
    private int port = 8080;

    @Override
    public Integer call() {
        try {
            AssinadorHttpServer server = new AssinadorHttpServer(port);
            server.start();

            System.out.println("Servidor HTTP do assinador iniciado com sucesso.");
            System.out.println("Porta: " + port);
            System.out.println("Health check: http://localhost:" + port + "/health");
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

            return 1;

        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.err.println("Servidor HTTP interrompido.");
            return 1;

        } catch (Exception e) {
            System.err.println("Erro inesperado ao iniciar servidor HTTP: " + e.getMessage());
            return 1;
        }
    }
}