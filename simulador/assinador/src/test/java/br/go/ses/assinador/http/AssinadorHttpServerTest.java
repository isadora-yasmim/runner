package br.go.ses.assinador.http;

import org.junit.jupiter.api.Test;

import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URI;

import static org.junit.jupiter.api.Assertions.assertEquals;

class AssinadorHttpServerTest {

    @Test
    void health_deveRetornarStatus200() throws Exception {
        AssinadorHttpServer server = new AssinadorHttpServer(9090);
        server.start();

        try {
            HttpURLConnection connection = (HttpURLConnection)
                    URI.create("http://localhost:9090/health").toURL().openConnection();

            connection.setRequestMethod("GET");

            assertEquals(200, connection.getResponseCode());
        } finally {
            server.stop();
        }
    }

    @Test
    void sign_comParametrosValidos_deveRetornarStatus200() throws Exception {
        AssinadorHttpServer server = new AssinadorHttpServer(9091);
        server.start();

        try {
            String json = "{\"document\":\"documento.txt\",\"tokenPin\":\"1234\"}";

            HttpURLConnection connection = (HttpURLConnection)
                    URI.create("http://localhost:9091/sign").toURL().openConnection();

            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setDoOutput(true);

            try (OutputStream os = connection.getOutputStream()) {
                os.write(json.getBytes());
            }

            assertEquals(200, connection.getResponseCode());
        } finally {
            server.stop();
        }
    }

    @Test
    void validate_comAssinaturaValida_deveRetornarStatus200() throws Exception {
        AssinadorHttpServer server = new AssinadorHttpServer(9092);
        server.start();

        try {
            String json = "{\"document\":\"documento.txt\",\"signature\":\"mock_hash_abc123_base64_encoded_signature_simulated\"}";

            HttpURLConnection connection = (HttpURLConnection)
                    URI.create("http://localhost:9092/validate").toURL().openConnection();

            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setDoOutput(true);

            try (OutputStream os = connection.getOutputStream()) {
                os.write(json.getBytes());
            }

            assertEquals(200, connection.getResponseCode());
        } finally {
            server.stop();
        }
    }

    @Test
    void sign_comMetodoGet_deveRetornarStatus405() throws Exception {
        AssinadorHttpServer server = new AssinadorHttpServer(9093);
        server.start();

        try {
            HttpURLConnection connection = (HttpURLConnection)
                    URI.create("http://localhost:9093/sign").toURL().openConnection();

            connection.setRequestMethod("GET");

            assertEquals(405, connection.getResponseCode());
        } finally {
            server.stop();
        }
    }
}