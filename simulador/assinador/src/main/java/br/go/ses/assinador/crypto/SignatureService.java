package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;
import br.go.ses.assinador.util.ParameterValidator;

import java.nio.charset.StandardCharsets;

/**
 * Recebedor das operacoes de assinatura. Concentra a <b>validacao de
 * parametros</b> (autoridade unica — criterio E3) e delega a operacao
 * criptografica a um {@link SignatureToken}, sem conhecer sua implementacao
 * concreta.
 *
 * <p>Ordem de tratamento de erro, do mais especifico ao mais geral:</p>
 * <ol>
 *   <li>parametro invalido -&gt; {@link IllegalArgumentException} (erro do usuario);</li>
 *   <li>dispositivo ausente/falho -&gt; {@link TokenException} (erro de sistema).</li>
 * </ol>
 * A validacao ocorre <b>antes</b> da checagem de presenca do dispositivo,
 * de modo que um erro do usuario nunca seja mascarado por um erro de sistema.
 */
public class SignatureService {

    private final SignatureToken token;

    public SignatureService(SignatureToken token) {
        this.token = token;
    }

    public SignatureData sign(String document, String pin) throws TokenException {
        ParameterValidator.validateSign(document, pin);
        requirePresentToken();
        return token.sign(
                document.getBytes(StandardCharsets.UTF_8),
                pin.toCharArray());
    }

    public boolean verify(String document, String signature) throws TokenException {
        ParameterValidator.validateVerify(document, signature);
        requirePresentToken();
        return token.verify(
                document.getBytes(StandardCharsets.UTF_8),
                signature);
    }

    /** Descricao do token em uso, para diagnostico. */
    public String tokenDescription() {
        return token.describe();
    }

    private void requirePresentToken() throws TokenException {
        if (!token.isPresent()) {
            throw new TokenException(
                    "Dispositivo de assinatura nao disponivel: " + token.describe()
                            + ". Utilize o modo simulado (sem --real) ou configure o dispositivo.");
        }
    }
}
