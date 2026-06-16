package br.go.ses.assinador.commands;

import br.go.ses.assinador.crypto.SignatureService;
import br.go.ses.assinador.crypto.TokenException;
import br.go.ses.assinador.crypto.TokenFactory;
import br.go.ses.assinador.model.SignatureData;
import br.go.ses.assinador.util.JsonUtil;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(
    name = "sign",
    description = "Cria uma assinatura digital (simulada por padrao)"
)
public class SignCommand implements Callable<Integer> {

    @Option(
        names = {"-d", "--document"},
        description = "Caminho ou conteudo do documento a ser assinado",
        required = true
    )
    private String document;

    @Option(
        names = {"-t", "--token-pin"},
        description = "PIN do dispositivo criptografico",
        required = true
    )
    private String tokenPin;

    @Option(
        names = {"--real"},
        description = "Usa dispositivo PKCS#11 real em vez do token simulado (padrao: simulado)"
    )
    private boolean real;

    @Override
    public Integer call() {
        SignatureService service = new SignatureService(TokenFactory.create(real));

        try {
            SignatureData signature = service.sign(document, tokenPin);
            JsonUtil.printSuccess("Assinatura criada com sucesso (Simulacao)", signature);
            return 0;

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
