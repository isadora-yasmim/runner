package br.go.ses.assinador.http;

import br.go.ses.assinador.model.ResponseOutput;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.Map;

public class AssinadorHttpServer {

    private final int port;
    private HttpServer server;
    private static final ObjectMapper mapper = new ObjectMapper();

    public AssinadorHttpServer(int port) {
        this.port = port;
    }

    public void start() throws IOException {
        server = HttpServer.create(new InetSocketAddress(port), 0);

        // =========================
        // HEALTH CHECK
        // =========================
        server.createContext("/health", exchange -> {
            try {
                if (!"GET".equals(exchange.getRequestMethod())) {
                    sendResponse(exchange, 405,
                        new ResponseOutput(false, "Metodo nao permitido", null));
                    return;
                }

                ResponseOutput response = new ResponseOutput(
                        true,
                        "Assinador HTTP ativo",
                        Map.of("status", "UP")
                );

                sendResponse(exchange, 200, response);

            } catch (Exception e) {
                sendResponse(exchange, 500,
                    new ResponseOutput(false, "Erro interno: " + e.getMessage(), null));
            }
        });

        server.setExecutor(null);
        server.start();
    }

    public void stop() {
        if (server != null) {
            server.stop(0);
        }
    }

    // =========================
    // MÉTODO UTILITÁRIO
    // =========================
    private void sendResponse(
            com.sun.net.httpserver.HttpExchange exchange,
            int statusCode,
            ResponseOutput response
    ) throws IOException {

        String json = mapper
                .writerWithDefaultPrettyPrinter()
                .writeValueAsString(response);

        byte[] bytes = json.getBytes();

        exchange.getResponseHeaders().add("Content-Type", "application/json");
        exchange.sendResponseHeaders(statusCode, bytes.length);

        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}