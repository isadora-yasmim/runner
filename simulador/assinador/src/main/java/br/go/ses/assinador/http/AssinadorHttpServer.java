package br.go.ses.assinador.http;

import br.go.ses.assinador.crypto.SignatureService;
import br.go.ses.assinador.crypto.SignatureToken;
import br.go.ses.assinador.crypto.TokenException;
import br.go.ses.assinador.crypto.TokenFactory;
import br.go.ses.assinador.model.ResponseOutput;
import br.go.ses.assinador.model.SignatureData;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Map;

public class AssinadorHttpServer {

    private final int port;
    private final SignatureService signatureService;
    private HttpServer server;
    private boolean testMode = false;

    private static final ObjectMapper mapper = new ObjectMapper();

    /**
     * Construtor de conveniencia: mantem compatibilidade com chamadas
     * existentes, usando o token simulado por padrao.
     */
    public AssinadorHttpServer(int port) {
        this(port, TokenFactory.create(false));
    }

    /**
     * Construtor principal: o transporte (HTTP) recebe o token a ser usado
     * pelo dominio, sem conhecer sua implementacao concreta.
     */
    public AssinadorHttpServer(int port, SignatureToken token) {
        this.port = port;
        this.signatureService = new SignatureService(token);
    }

    /**
     * Ativa o modo de teste: o endpoint /stop encerra apenas o HttpServer,
     * sem chamar System.exit(), evitando matar a JVM do Surefire.
     */
    public AssinadorHttpServer withTestMode() {
        this.testMode = true;
        return this;
    }

    /** Porta real atribuída pelo SO (útil com porta 0 em testes). */
    public int getPort() {
        if (server == null) {
            throw new IllegalStateException("Servidor nao iniciado.");
        }
        return server.getAddress().getPort();
    }

    public void start() throws IOException {
        server = HttpServer.create(new InetSocketAddress(port), 0);

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
                        Map.of("status", "UP", "token", signatureService.tokenDescription())
                );

                sendResponse(exchange, 200, response);

            } catch (IOException e) {
                sendResponse(exchange, 500,
                        new ResponseOutput(false, "Erro interno: " + e.getMessage(), null));
            }
        });

        server.createContext("/sign", exchange -> {
            try {
                if (!"POST".equals(exchange.getRequestMethod())) {
                    sendResponse(exchange, 405,
                            new ResponseOutput(false, "Metodo nao permitido", null));
                    return;
                }

                JsonNode body = readJsonBody(exchange);
                String document = getTextField(body, "document");
                String tokenPin = getTextField(body, "tokenPin");

                SignatureData signature = signatureService.sign(document, tokenPin);

                sendResponse(exchange, 200,
                        new ResponseOutput(true, "Assinatura criada com sucesso (Simulacao)", signature));

            } catch (IllegalArgumentException e) {
                // Erro do usuario: parametro invalido.
                sendResponse(exchange, 400,
                        new ResponseOutput(false, e.getMessage(), null));

            } catch (TokenException e) {
                // Erro de sistema: dispositivo ausente ou falho.
                sendResponse(exchange, 503,
                        new ResponseOutput(false, e.getMessage(), null));

            } catch (IOException e) {
                sendResponse(exchange, 500,
                        new ResponseOutput(false, "Erro interno: " + e.getMessage(), null));
            }
        });

        server.createContext("/validate", exchange -> {
            try {
                if (!"POST".equals(exchange.getRequestMethod())) {
                    sendResponse(exchange, 405,
                            new ResponseOutput(false, "Metodo nao permitido", null));
                    return;
                }

                JsonNode body = readJsonBody(exchange);
                String document = getTextField(body, "document");
                String signature = getTextField(body, "signature");

                boolean isValid = signatureService.verify(document, signature);

                if (isValid) {
                    sendResponse(exchange, 200,
                            new ResponseOutput(true, "Assinatura valida.", true));
                } else {
                    sendResponse(exchange, 200,
                            new ResponseOutput(false, "Assinatura invalida ou corrompida.", false));
                }

            } catch (IllegalArgumentException e) {
                sendResponse(exchange, 400,
                        new ResponseOutput(false, e.getMessage(), null));

            } catch (TokenException e) {
                sendResponse(exchange, 503,
                        new ResponseOutput(false, e.getMessage(), null));

            } catch (IOException e) {
                sendResponse(exchange, 500,
                        new ResponseOutput(false, "Erro interno: " + e.getMessage(), null));
            }
        });

        server.createContext("/stop", exchange -> {
            try {
                if (!"POST".equals(exchange.getRequestMethod())) {
                    sendResponse(exchange, 405,
                            new ResponseOutput(false, "Metodo nao permitido", null));
                    return;
                }

                sendResponse(exchange, 200,
                        new ResponseOutput(true, "Servidor encerrado com sucesso.", null));

                Thread shutdownThread = new Thread(() -> {
                    try {
                        Thread.sleep(500);
                        server.stop(0);
                        if (!testMode) {
                            System.exit(0);
                        }
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        System.err.println("Encerramento interrompido: " + e.getMessage());
                        if (!testMode) {
                            System.exit(1);
                        }
                    }
                });

                shutdownThread.setDaemon(false);
                shutdownThread.start();

            } catch (IOException e) {
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

    private JsonNode readJsonBody(HttpExchange exchange) throws IOException {
        String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);

        if (body.trim().isEmpty()) {
            throw new IllegalArgumentException("Corpo da requisicao nao pode estar vazio.");
        }

        return mapper.readTree(body);
    }

    private String getTextField(JsonNode body, String fieldName) {
        if (body == null || !body.has(fieldName) || body.get(fieldName).isNull()) {
            return null;
        }

        return body.get(fieldName).asText();
    }

    private void sendResponse(
            HttpExchange exchange,
            int statusCode,
            ResponseOutput response
    ) throws IOException {

        String json = mapper
                .writerWithDefaultPrettyPrinter()
                .writeValueAsString(response);

        byte[] bytes = json.getBytes(StandardCharsets.UTF_8);

        exchange.getResponseHeaders().add("Content-Type", "application/json");
        exchange.sendResponseHeaders(statusCode, bytes.length);

        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}
