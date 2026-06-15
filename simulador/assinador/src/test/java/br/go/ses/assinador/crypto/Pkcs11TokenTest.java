package br.go.ses.assinador.crypto;

import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

class Pkcs11TokenTest {

    private final Pkcs11Token token = new Pkcs11Token();

    @Test
    void isPresent_semConfiguracao_deveSerFalso() {
        // Esqueleto: deteccao real sera implementada na Issue #47.
        assertFalse(token.isPresent());
    }

    @Test
    void sign_semConfiguracao_deveLancarTokenException() {
        TokenException ex = assertThrows(TokenException.class,
                () -> token.sign("doc".getBytes(StandardCharsets.UTF_8), "1234".toCharArray()));

        // Mensagem orientativa, nao stack trace.
        assertTrue(ex.getMessage().contains("PKCS#11"));
    }

    @Test
    void verify_semConfiguracao_deveLancarTokenException() {
        assertThrows(TokenException.class,
                () -> token.verify("doc".getBytes(StandardCharsets.UTF_8), "qualquer"));
    }

    @Test
    void describe_deveIdentificarTokenPkcs11() {
        assertTrue(token.describe().contains("PKCS#11"));
    }
}
