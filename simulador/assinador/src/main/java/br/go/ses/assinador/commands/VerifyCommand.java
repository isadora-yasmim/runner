package br.go.ses.assinador.commands;

import br.go.ses.assinador.util.JsonUtil;
import br.go.ses.assinador.util.ParameterValidator;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(name = "verify", description = "Simula a validacao de uma assinatura digital")
public class VerifyCommand implements Callable<Integer> {

    @Option(names = {"-d", "--document"}, description = "Documento original", required = true)
    private String document;

    @Option(names = {"-s", "--signature"}, description = "Hash ou conteudo da assinatura", required = true)
    private String signature;

    @Override
    public Integer call() {

        try {
            ParameterValidator.validateVerify(document, signature);
        } catch (IllegalArgumentException e) {
            JsonUtil.printError(e.getMessage());
            return 1;
        }

        if ("mock_hash_abc123_base64_encoded_signature_simulated".equals(signature)) {
            JsonUtil.printSuccess("Assinatura valida.", true);
            return 0;
        } else {
            JsonUtil.printError("Assinatura invalida ou corrompida.");
            return 1;
        }
    }
}