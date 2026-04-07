package br.go.ses.assinador.commands;
import br.go.ses.assinador.util.JsonUtil;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(name = "verify", description = "Simula a validação de uma assinatura digital")
public class VerifyCommand implements Callable<Integer> {

    @Option(names = {"-d", "--document"}, description = "Documento original", required = true)
    private String document;

    @Option(names = {"-s", "--signature"}, description = "Hash ou conteúdo da assinatura", required = true)
    private String signature;

    @Override
    public Integer call() {
        // Validação rigorosa dos parâmetros
        if (document.trim().isEmpty() || signature.trim().isEmpty()) {
            JsonUtil.printError("Documento e Assinatura são obrigatórios para a validação.");
            return 1;
        }

        // Lógica simples de simulação: se a assinatura for uma string específica, é válida
        if ("mock_hash_abc123_base64_encoded_signature_simulated".equals(signature)) {
            JsonUtil.printSuccess("Assinatura válida.", true);
            return 0;
        } else {
            JsonUtil.printError("Assinatura inválida ou corrompida.");
            return 1;
        }
    }
}