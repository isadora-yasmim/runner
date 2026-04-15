package br.go.ses.assinador.util;

public class ParameterValidator {

    public static void validateSign(String document, String tokenPin) {
        if (document == null || document.trim().isEmpty()) {
            throw new IllegalArgumentException("O documento nao pode estar vazio.");
        }

        if (tokenPin == null || tokenPin.trim().isEmpty()) {
            throw new IllegalArgumentException("O PIN do token eh obrigatorio.");
        }

        if (tokenPin.trim().length() < 4) {
            throw new IllegalArgumentException("O PIN do token deve conter pelo menos 4 caracteres.");
        }
    }

    public static void validateVerify(String document, String signature) {
        if (document == null || document.trim().isEmpty()) {
            throw new IllegalArgumentException("O documento eh obrigatorio para a validacao.");
        }

        if (signature == null || signature.trim().isEmpty()) {
            throw new IllegalArgumentException("A assinatura eh obrigatoria para a validacao.");
        }
    }
}