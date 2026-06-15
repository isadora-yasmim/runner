package br.go.ses.assinador.commands;

import br.go.ses.assinador.crypto.SignatureService;
import br.go.ses.assinador.crypto.TokenException;
import br.go.ses.assinador.crypto.TokenFactory;
import br.go.ses.assinador.util.JsonUtil;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(
    name = "verify",
    description = "Valida uma assinatura digital (simulada por padrao)"
)
public class VerifyCommand implements Callable<Integer> {

    @Option(
        names = {"-d", "--document"},
        description = "Documento original",
        required = true
    )
    private String document;

    @Option(
        names = {"-s", "--signature"},
        description = "Hash ou conteudo da assinatura",
        required = true
    )
    private String signature;

    @Option(
        names = {"--real"},
        description = "Usa dispositivo PKCS#11 real em vez do token simulado (padrao: simulado)"
    )
    private boolean real;

    @Override
    public Integer call() {
        SignatureService service = new SignatureService(TokenFactory.create(real));

        try {
            boolean valid = service.verify(document, signature);
            if (valid) {
                JsonUtil.printSuccess("Assinatura valida.", true);
                return 0;
            }
            // Assinatura invalida e um resultado do dominio, nao um erro de sistema.
            JsonUtil.printError("Assinatura invalida ou corrompida.");
            return 1;

        } catch (IllegalArgumentException e) {
            // Erro do usuario: parametro invalido.
            JsonUtil.printError(e.getMessage());
            return 1;

        } catch (TokenException e) {
            // Erro de sistema: dispositivo ausente ou falho.
            JsonUtil.printError(e.getMessage());
            return 2;
        }
    }
}
