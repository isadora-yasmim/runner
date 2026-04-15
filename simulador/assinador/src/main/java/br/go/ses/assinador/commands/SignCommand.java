package br.go.ses.assinador.commands;

import br.go.ses.assinador.model.SignatureData;
import br.go.ses.assinador.util.JsonUtil;
import br.go.ses.assinador.util.ParameterValidator;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(name = "sign", description = "Simula a criacao de uma assinatura digital")
public class SignCommand implements Callable<Integer> {

    @Option(names = {"-d", "--document"}, description = "Caminho ou conteudo do documento a ser assinado", required = true)
    private String document;

    @Option(names = {"-t", "--token-pin"}, description = "PIN do dispositivo criptografico", required = true)
    private String tokenPin;

    @Override
    public Integer call() {
        try {
            ParameterValidator.validateSign(document, tokenPin);
        } catch (IllegalArgumentException e) {
            JsonUtil.printError(e.getMessage());
            return 1;
        }

        SignatureData mockSignature = new SignatureData(
            "mock_hash_abc123_base64_encoded_signature_simulated",
            "SHA256withRSA"
        );

        JsonUtil.printSuccess("Assinatura criada com sucesso (Simulacao)", mockSignature);
        return 0;
    }
}