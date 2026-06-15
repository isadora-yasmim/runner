package br.go.ses.assinador.crypto;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TokenFactoryTest {

    @Test
    void create_modoSimulado_deveRetornarSimulatedToken() {
        SignatureToken token = TokenFactory.create(TokenFactory.Mode.SIMULATED);
        assertInstanceOf(SimulatedToken.class, token);
    }

    @Test
    void create_modoReal_deveRetornarPkcs11Token() {
        SignatureToken token = TokenFactory.create(TokenFactory.Mode.REAL);
        assertInstanceOf(Pkcs11Token.class, token);
    }

    @Test
    void create_booleanoFalse_deveRetornarSimulatedToken() {
        assertInstanceOf(SimulatedToken.class, TokenFactory.create(false));
    }

    @Test
    void create_booleanoTrue_deveRetornarPkcs11Token() {
        assertInstanceOf(Pkcs11Token.class, TokenFactory.create(true));
    }
}
