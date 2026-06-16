package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;

/**
 * Token simulado: in-process, deterministico, sem dependencia nativa.
 *
 * <p>Implementa o mesmo contrato {@link SignatureToken} de um dispositivo
 * real e, do ponto de vista do {@link SignatureService}, e indistinguivel
 * dele — "engana" o recebedor se passando por um token/smartcard.</p>
 *
 * <p>E a <b>unica fonte da verdade</b> do material simulado: o hash e o
 * algoritmo deixaram de estar duplicados em {@code SignCommand} e
 * {@code AssinadorHttpServer}.</p>
 */
public class SimulatedToken implements SignatureToken {

    /** Assinatura simulada conhecida; reconhecida como valida na verificacao. */
    public static final String MOCK_SIGNATURE =
            "mock_hash_abc123_base64_encoded_signature_simulated";

    private static final String ALGORITHM = "SHA256withRSA";

    @Override
    public SignatureData sign(byte[] content, char[] pin) {
        // Simulacao: retorna sempre a assinatura pre-construida.
        return new SignatureData(MOCK_SIGNATURE, ALGORITHM);
    }

    @Override
    public boolean verify(byte[] content, String signature) {
        return MOCK_SIGNATURE.equals(signature);
    }

    @Override
    public boolean isPresent() {
        // O simulador esta sempre disponivel.
        return true;
    }

    @Override
    public String describe() {
        return "Token simulado (in-process, deterministico)";
    }
}
