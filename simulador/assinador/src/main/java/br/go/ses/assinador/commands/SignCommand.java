package br.go.ses.assinador.commands;
import br.go.ses.assinador.model.SignatureData;
import br.go.ses.assinador.util.JsonUtil;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(name = "sign", description = "Simula a criação de uma assinatura digital")
public class SignCommand implements Callable<Integer> {

    @Option(names = {"-d", "--document"}, description = "Caminho ou conteúdo do documento a ser assinado", required = true)
    private String document;

    @Option(names = {"-t", "--token-pin"}, description = "PIN do dispositivo criptográfico", required = true)
    private String tokenPin;

    @Override
    public Integer call() {
        // Validação rigorosa dos parâmetros (Simulação)
        if (document.trim().isEmpty()) {
            JsonUtil.printError("O documento não pode estar vazio.");
            return 1; // Código de erro genérico
        }

        if (tokenPin.length() < 4) {
            JsonUtil.printError("O PIN do token deve conter pelo menos 4 caracteres.");
            return 1;
        }

        // Simulação de criação de assinatura retornando resposta pré-construída
        SignatureData mockSignature = new SignatureData(
            "mock_hash_abc123_base64_encoded_signature_simulated",
            "SHA256withRSA"
        );

        JsonUtil.printSuccess("Assinatura criada com sucesso (Simulação)", mockSignature);
        return 0; // Sucesso
    }
}