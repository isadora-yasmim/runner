package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;

/**
 * Token real via PKCS#11 (provider {@code SunPKCS11}).
 *
 * <p><b>Esqueleto</b>: a implementacao efetiva — incluindo testes de
 * integracao com SoftHSM2 — sera feita na Issue #47 (US-02.5). Por ora:</p>
 * <ul>
 *   <li>{@link #isPresent()} retorna {@code false} (modulo nao configurado);</li>
 *   <li>{@link #sign} e {@link #verify} lancam {@link TokenException} com
 *       mensagem orientativa, em vez de stack trace.</li>
 * </ul>
 *
 * <p>Como respeita o contrato {@link SignatureToken}, plugar a implementacao
 * real na Issue #47 nao exigira mudanca alguma no {@link SignatureService},
 * no CLI ou no servidor HTTP.</p>
 */
public class Pkcs11Token implements SignatureToken {

    private static final String NOT_CONFIGURED =
            "Operacao via PKCS#11 ainda nao configurada. "
            + "Configure o modulo PKCS#11 / SoftHSM2 (ver Issue #47) "
            + "ou utilize o modo simulado (sem a flag --real).";

    @Override
    public SignatureData sign(byte[] content, char[] pin) throws TokenException {
        throw new TokenException(NOT_CONFIGURED);
    }

    @Override
    public boolean verify(byte[] content, String signature) throws TokenException {
        throw new TokenException(NOT_CONFIGURED);
    }

    @Override
    public boolean isPresent() {
        // Deteccao real do modulo PKCS#11 sera implementada na Issue #47.
        return false;
    }

    @Override
    public String describe() {
        return "Token PKCS#11 (nao configurado — ver Issue #47)";
    }
}
