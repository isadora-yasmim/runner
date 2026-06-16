package br.go.ses.assinador.crypto;

import br.go.ses.assinador.model.SignatureData;

/**
 * Porta (contrato) de um dispositivo de assinatura: token, smartcard ou
 * simulador. Quem recebe o pedido de assinatura ({@link SignatureService})
 * depende apenas desta abstracao e nao conhece a implementacao concreta.
 *
 * <p>Por isso um token simulado ({@link SimulatedToken}) e indistinguivel,
 * do ponto de vista do recebedor, de um token real via PKCS#11
 * ({@link Pkcs11Token}).</p>
 */
public interface SignatureToken {

    /**
     * Cria uma assinatura sobre o conteudo informado.
     *
     * @param content conteudo a ser assinado (bytes)
     * @param pin     PIN do dispositivo criptografico
     * @return os dados da assinatura produzida
     * @throws TokenException quando o dispositivo falha ou nao esta disponivel
     *                        (erro de sistema, distinto de erro de parametro)
     */
    SignatureData sign(byte[] content, char[] pin) throws TokenException;

    /**
     * Verifica se a assinatura informada corresponde ao conteudo.
     *
     * @param content   conteudo original (bytes)
     * @param signature assinatura a ser verificada
     * @return {@code true} se a assinatura for valida
     * @throws TokenException quando o dispositivo falha ou nao esta disponivel
     */
    boolean verify(byte[] content, String signature) throws TokenException;

    /**
     * Indica se o dispositivo esta presente e pronto para uso.
     * Health check do token, anterior a qualquer operacao.
     *
     * @return {@code true} se o dispositivo esta disponivel
     */
    boolean isPresent();

    /**
     * Descricao legivel do dispositivo, para diagnostico e logs.
     *
     * @return texto curto identificando o token
     */
    String describe();
}
