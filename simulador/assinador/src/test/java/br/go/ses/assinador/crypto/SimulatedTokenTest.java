package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

class SimulatedTokenTest {

    private final SimulatedToken token = new SimulatedToken();

    @Test
    void isPresent_deveSerSempreVerdadeiro() {
        assertTrue(token.isPresent());
    }

    @Test
    void sign_deveRetornarHashEAlgoritmoSimulados() throws Exception {
        SignatureData data = token.sign("doc".getBytes(StandardCharsets.UTF_8), "1234".toCharArray());

        assertEquals(SimulatedToken.MOCK_SIGNATURE, data.getSignatureHash());
        assertEquals("SHA256withRSA", data.getAlgorithm());
    }

    @Test
    void verify_comAssinaturaSimulada_deveSerValida() throws Exception {
        boolean valid = token.verify(
                "doc".getBytes(StandardCharsets.UTF_8),
                SimulatedToken.MOCK_SIGNATURE);

        assertTrue(valid);
    }

    @Test
    void verify_comAssinaturaDiferente_deveSerInvalida() throws Exception {
        boolean valid = token.verify(
                "doc".getBytes(StandardCharsets.UTF_8),
                "assinatura_errada");

        assertFalse(valid);
    }

    @Test
    void describe_deveIdentificarTokenSimulado() {
        assertTrue(token.describe().toLowerCase().contains("simulado"));
    }
}
