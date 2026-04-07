package br.go.ses.assinador.model;
public class SignatureData {
    private String signatureHash;
    private String algorithm;

    public SignatureData(String signatureHash, String algorithm) {
        this.signatureHash = signatureHash;
        this.algorithm = algorithm;
    }

    public String getSignatureHash() { return signatureHash; }
    public String getAlgorithm() { return algorithm; }
}