package br.go.ses.assinador.http;

import br.go.ses.assinador.model.ResponseOutput;
import br.go.ses.assinador.model.SignatureData;
import br.go.ses.assinador.util.ParameterValidator;

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
    private HttpServer server;
    private static final ObjectMapper mapper = new ObjectMapper();

    private static final String MOCK_SIGNATURE =
            "mock_hash_abc123_base64_encoded_signature_simulated";

    public AssinadorHttpServer(int port) {
        this.port = port;
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
                        Map.of("status", "UP")
                );

                sendResponse(exchange, 200, response);

            } catch (Exception e) {
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

                ParameterValidator.validateSign(document, tokenPin);

                SignatureData signatureData = new SignatureData(
                        MOCK_SIGNATURE,
                        "SHA256withRSA"
                );

                sendResponse(exchange, 200,
                        new ResponseOutput(true, "Assinatura criada com sucesso (Simulacao)", signatureData));

            } catch (IllegalArgumentException e) {
                sendResponse(exchange, 400,
                        new ResponseOutput(false, e.getMessage(), null));

            } catch (Exception e) {
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

                ParameterValidator.validateVerify(document, signature);

                boolean isValid = MOCK_SIGNATURE.equals(signature);

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

            } catch (Exception e) {
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
                        System.exit(0);
                    } catch (Exception e) {
                        System.err.println("Erro ao encerrar servidor: " + e.getMessage());
                        System.exit(1);
                    }
                });

                shutdownThread.setDaemon(false);
                shutdownThread.start();

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