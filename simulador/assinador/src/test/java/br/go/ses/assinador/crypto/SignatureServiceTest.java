package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SignatureServiceTest {

    private final SignatureService simulado =
            new SignatureService(new SimulatedToken());

    // --- caminho feliz com token simulado ---

    @Test
    void sign_parametrosValidos_deveRetornarAssinatura() throws Exception {
        SignatureData data = simulado.sign("doc.pdf", "1234");

        assertEquals(SimulatedToken.MOCK_SIGNATURE, data.getSignatureHash());
        assertEquals("SHA256withRSA", data.getAlgorithm());
    }

    @Test
    void verify_assinaturaValida_deveRetornarTrue() throws Exception {
        assertTrue(simulado.verify("doc.pdf", SimulatedToken.MOCK_SIGNATURE));
    }

    @Test
    void verify_assinaturaInvalida_deveRetornarFalse() throws Exception {
        assertFalse(simulado.verify("doc.pdf", "hash_errado"));
    }

    // --- erro do usuario: validacao tem precedencia ---

    @Test
    void sign_pinCurto_deveLancarIllegalArgument() {
        assertThrows(IllegalArgumentException.class,
                () -> simulado.sign("doc.pdf", "12"));
    }

    @Test
    void sign_documentoVazio_deveLancarIllegalArgument() {
        assertThrows(IllegalArgumentException.class,
                () -> simulado.sign("   ", "1234"));
    }

    @Test
    void verify_assinaturaVazia_deveLancarIllegalArgument() {
        assertThrows(IllegalArgumentException.class,
                () -> simulado.verify("doc.pdf", ""));
    }

    // --- erro de sistema: dispositivo ausente ---

    @Test
    void sign_dispositivoAusente_deveLancarTokenException() {
        SignatureService real = new SignatureService(new Pkcs11Token());

        // Parametros validos: o erro vem do dispositivo, nao do usuario.
        assertThrows(TokenException.class,
                () -> real.sign("doc.pdf", "1234"));
    }

    @Test
    void sign_parametroInvalidoTemPrecedenciaSobreDispositivoAusente() {
        SignatureService real = new SignatureService(new Pkcs11Token());

        // Mesmo com dispositivo ausente, parametro invalido deve falhar primeiro
        // (erro do usuario antes de erro de sistema).
        assertThrows(IllegalArgumentException.class,
                () -> real.sign("doc.pdf", "12"));
    }

    @Test
    void tokenDescription_deveRefletirTokenEmUso() {
        assertTrue(simulado.tokenDescription().toLowerCase().contains("simulado"));
    }
}
