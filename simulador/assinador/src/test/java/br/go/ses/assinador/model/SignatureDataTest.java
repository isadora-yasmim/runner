// SignatureDataTest.java
package br.go.ses.assinador.model;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class SignatureDataTest {

    @Test
    void construtor_devePreencherCamposCorretamente() {
        SignatureData data = new SignatureData("hash123", "SHA256withRSA");
        assertEquals("hash123", data.getSignatureHash());
        assertEquals("SHA256withRSA", data.getAlgorithm());
    }
}