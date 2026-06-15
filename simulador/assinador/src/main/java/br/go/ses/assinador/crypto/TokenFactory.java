package br.go.ses.assinador.crypto;

/**
 * Seleciona a implementacao de {@link SignatureToken} a ser usada.
 *
 * <p>Modo <b>simulado</b> e o padrao; o modo <b>real</b> (PKCS#11) deve ser
 * ativado explicitamente — alinhado ao principio do projeto de que a
 * simulacao e o caminho padrao e o uso de dispositivo real e opt-in.</p>
 */
public final class TokenFactory {

    public enum Mode {
        SIMULATED,
        REAL
    }

    private TokenFactory() {
        // utilitario; nao instanciavel
    }

    public static SignatureToken create(Mode mode) {
        return switch (mode) {
            case REAL -> new Pkcs11Token();
            case SIMULATED -> new SimulatedToken();
        };
    }

    /**
     * Conveniencia para a CLI: {@code true} seleciona o token real (PKCS#11),
     * {@code false} (padrao) seleciona o token simulado.
     */
    public static SignatureToken create(boolean real) {
        return create(real ? Mode.REAL : Mode.SIMULATED);
    }
}
