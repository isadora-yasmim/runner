package br.go.ses.assinador.crypto;

/**
 * Erro de <b>sistema/dispositivo</b> durante uma operacao de assinatura:
 * token ausente, modulo PKCS#11 nao configurado, falha de comunicacao com
 * o hardware, etc.
 *
 * <p>Deliberadamente distinta de {@link IllegalArgumentException}, que
 * representa <b>erro do usuario</b> (parametro invalido). Essa separacao
 * permite que a interface (CLI/HTTP) traduza cada categoria em codigos de
 * saida e status HTTP diferentes:</p>
 * <ul>
 *   <li>{@code IllegalArgumentException} -&gt; exit code 1 / HTTP 400</li>
 *   <li>{@code TokenException}          -&gt; exit code 2 / HTTP 503</li>
 * </ul>
 *
 * <p>E uma <i>checked exception</i> de proposito: forca o tratamento
 * explicito em quem invoca a porta, evitando que falhas de dispositivo
 * sejam silenciosamente engolidas.</p>
 */
public class TokenException extends Exception {

    public TokenException(String message) {
        super(message);
    }

    public TokenException(String message, Throwable cause) {
        super(message, cause);
    }
}
