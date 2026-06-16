package br.go.ses.assinador.http;

import br.go.ses.assinador.crypto.TokenFactory;
import org.junit.jupiter.api.Test;

import java.net.HttpURLConnection;
import java.net.URI;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Valida o auto-shutdown por inatividade (US-01.9 / criterio E2),
 * em especial que o timer reinicia a cada requisicao recebida.
 * Usa withTestMode() para que o encerramento por inatividade nao
 * chame System.exit() e mate a JVM do Surefire.
 */
class AssinadorHttpServerTimeoutTest {

    private int hitHealth(int port) throws Exception {
        HttpURLConnection conn = (HttpURLConnection)
                URI.create("http://localhost:" + port + "/health").toURL().openConnection();
        conn.setRequestMethod("GET");
        int code = conn.getResponseCode();
        conn.disconnect();
        return code;
    }

    @Test
    void timeoutZero_naoAgendaShutdown() throws Exception {
        AssinadorHttpServer server =
                new AssinadorHttpServer(0, TokenFactory.create(false), 0L).withTestMode();
        server.start();
        int port = server.getPort();

        try {
            // janela mais que suficiente; com timeout 0 o servidor segue vivo
            Thread.sleep(1200);
            assertEquals(200, hitHealth(port));
        } finally {
            server.stop();
        }
    }

    @Test
    void requisicoesReiniciamOTimer_eEncerraQuandoTrafegoPara() throws Exception {
        // janela curta (600ms) via construtor package-private em ms
        AssinadorHttpServer server = new AssinadorHttpServer(
                0, TokenFactory.create(false), 600L, TimeUnit.MILLISECONDS).withTestMode();
        server.start();
        int port = server.getPort();

        try {
            // mantem vivo: 4 requisicoes espacadas 300ms (< janela de 600ms)
            // tempo total ~1200ms, bem alem da janela original de 600ms
            for (int i = 0; i < 4; i++) {
                assertEquals(200, hitHealth(port), "Servidor deveria responder enquanto ativo");
                Thread.sleep(300);
            }

            // agora paramos o trafego e esperamos alem da janela:
            // o servidor deve ter encerrado por inatividade
            Thread.sleep(1500);

            assertThrows(Exception.class, () -> hitHealth(port),
                    "Servidor deveria ter encerrado apos a inatividade");
        } finally {
            server.stop();
        }
    }
}