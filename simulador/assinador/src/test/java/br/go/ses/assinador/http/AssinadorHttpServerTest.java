package br.go.ses.assinador.http;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Testes de integração do AssinadorHttpServer.
 *
 * Todos os servidores são criados com withTestMode(), que faz o endpoint
 * /stop encerrar apenas o HttpServer sem chamar System.exit(), evitando
 * que a JVM forkeada do Surefire seja morta prematuramente.
 *
 * Usa porta 0 para que o SO escolha uma porta livre em cada teste,
 * eliminando conflitos em CI e execuções paralelas (Issue 43).
 */
class AssinadorHttpServerTest {

    private static final ObjectMapper mapper = new ObjectMapper();

    private AssinadorHttpServer server;
    private int port;

    @BeforeEach
    void setUp() throws Exception {
        server = new AssinadorHttpServer(0).withTestMode();
        server.start();
        port = server.getPort();
    }

    @AfterEach
    void tearDown() {
        server.stop();
    }

    // ------------------------------------------------------------------ /health

    @Test
    void health_GET_deveRetornarStatus200() throws Exception {
        HttpURLConnection conn = get("/health");
        assertEquals(200, conn.getResponseCode());
    }

    @Test
    void health_GET_deveRetornarSuccessTrueEStatusUP() throws Exception {
        HttpURLConnection conn = get("/health");
        JsonNode body = readBody(conn);
        assertTrue(body.get("success").asBoolean());
        assertEquals("UP", body.get("data").get("status").asText());
    }

    @Test
    void health_POST_deveRetornarStatus405() throws Exception {
        HttpURLConnection conn = post("/health", "{}");
        assertEquals(405, conn.getResponseCode());
    }

    // ------------------------------------------------------------------ /sign

    @Test
    void sign_POST_parametrosValidos_deveRetornarStatus200() throws Exception {
        String json = "{\"document\":\"documento.txt\",\"tokenPin\":\"1234\"}";
        HttpURLConnection conn = post("/sign", json);
        assertEquals(200, conn.getResponseCode());
    }

    @Test
    void sign_POST_parametrosValidos_deveRetornarHashEAlgoritmo() throws Exception {
        String json = "{\"document\":\"documento.txt\",\"tokenPin\":\"1234\"}";
        HttpURLConnection conn = post("/sign", json);
        JsonNode body = readBody(conn);
        assertTrue(body.get("success").asBoolean());
        assertEquals(
            "mock_hash_abc123_base64_encoded_signature_simulated",
            body.get("data").get("signatureHash").asText()
        );
        assertEquals("SHA256withRSA", body.get("data").get("algorithm").asText());
    }

    @Test
    void sign_POST_documentoAusente_deveRetornarStatus400() throws Exception {
        String json = "{\"tokenPin\":\"1234\"}";
        HttpURLConnection conn = post("/sign", json);
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void sign_POST_pinAusente_deveRetornarStatus400() throws Exception {
        String json = "{\"document\":\"doc.pdf\"}";
        HttpURLConnection conn = post("/sign", json);
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void sign_POST_pinMuitoCurto_deveRetornarStatus400() throws Exception {
        String json = "{\"document\":\"doc.pdf\",\"tokenPin\":\"12\"}";
        HttpURLConnection conn = post("/sign", json);
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void sign_POST_corpoVazio_deveRetornarStatus400() throws Exception {
        HttpURLConnection conn = post("/sign", "");
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void sign_POST_documentoComEspacos_deveRetornarStatus200() throws Exception {
        // Critério E1: argumentos com espaços devem ser preservados
        String json = "{\"document\":\"meu documento 2024.pdf\",\"tokenPin\":\"1234\"}";
        HttpURLConnection conn = post("/sign", json);
        assertEquals(200, conn.getResponseCode());
    }

    @Test
    void sign_GET_deveRetornarStatus405() throws Exception {
        HttpURLConnection conn = get("/sign");
        assertEquals(405, conn.getResponseCode());
    }

    // ------------------------------------------------------------------ /validate

    @Test
    void validate_POST_assinaturaValida_deveRetornarStatus200ESuccessTrue() throws Exception {
        String json = "{\"document\":\"doc.pdf\","
                + "\"signature\":\"mock_hash_abc123_base64_encoded_signature_simulated\"}";
        HttpURLConnection conn = post("/validate", json);
        assertEquals(200, conn.getResponseCode());
        JsonNode body = readBody(conn);
        assertTrue(body.get("success").asBoolean());
    }

    @Test
    void validate_POST_assinaturaInvalida_deveRetornarStatus200ESuccessFalse() throws Exception {
        String json = "{\"document\":\"doc.pdf\",\"signature\":\"hash_errado\"}";
        HttpURLConnection conn = post("/validate", json);
        assertEquals(200, conn.getResponseCode());
        JsonNode body = readBody(conn);
        assertFalse(body.get("success").asBoolean());
    }

    @Test
    void validate_POST_documentoAusente_deveRetornarStatus400() throws Exception {
        String json = "{\"signature\":\"qualquer\"}";
        HttpURLConnection conn = post("/validate", json);
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void validate_POST_assinaturaAusente_deveRetornarStatus400() throws Exception {
        String json = "{\"document\":\"doc.pdf\"}";
        HttpURLConnection conn = post("/validate", json);
        assertEquals(400, conn.getResponseCode());
    }

    @Test
    void validate_GET_deveRetornarStatus405() throws Exception {
        HttpURLConnection conn = get("/validate");
        assertEquals(405, conn.getResponseCode());
    }

    // ------------------------------------------------------------------ /stop

    @Test
    void stop_GET_deveRetornarStatus405() throws Exception {
        HttpURLConnection conn = get("/stop");
        assertEquals(405, conn.getResponseCode());
    }

    @Test
    void stop_POST_deveRetornarStatus200ESuccessTrue() throws Exception {
        // withTestMode() garante que /stop encerra só o HttpServer, sem System.exit()
        AssinadorHttpServer stopServer = new AssinadorHttpServer(0).withTestMode();
        stopServer.start();
        int stopPort = stopServer.getPort();

        try {
            HttpURLConnection conn = postOnPort(stopPort, "/stop", "{}");
            int status = conn.getResponseCode();
            JsonNode body = readBody(conn);

            assertEquals(200, status);
            assertTrue(body.get("success").asBoolean());
        } finally {
            stopServer.stop();
        }
    }

    // ------------------------------------------------------------------ helpers

    private HttpURLConnection get(String path) throws Exception {
        HttpURLConnection conn = openConnection(port, path);
        conn.setRequestMethod("GET");
        return conn;
    }

    private HttpURLConnection post(String path, String jsonBody) throws Exception {
        return postOnPort(port, path, jsonBody);
    }

    private HttpURLConnection postOnPort(int targetPort, String path, String jsonBody) throws Exception {
        HttpURLConnection conn = openConnection(targetPort, path);
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setDoOutput(true);
        try (OutputStream os = conn.getOutputStream()) {
            os.write(jsonBody.getBytes(StandardCharsets.UTF_8));
        }
        return conn;
    }

    private HttpURLConnection openConnection(int targetPort, String path) throws Exception {
        return (HttpURLConnection)
                URI.create("http://localhost:" + targetPort + path).toURL().openConnection();
    }

    private JsonNode readBody(HttpURLConnection conn) throws Exception {
        InputStream stream;
        try {
            stream = conn.getInputStream();
        } catch (Exception e) {
            stream = conn.getErrorStream();
        }
        if (stream == null) {
            return mapper.createObjectNode();
        }
        return mapper.readTree(stream.readAllBytes());
    }
}