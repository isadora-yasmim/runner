package br.go.ses.assinador;

import br.go.ses.assinador.commands.SignCommand;
import br.go.ses.assinador.commands.VerifyCommand;
import picocli.CommandLine;
import picocli.CommandLine.Command;

import java.util.concurrent.Callable;

@Command(
    name = "assinador",
    mixinStandardHelpOptions = true,
    version = "1.0",
    description = "Assinador Simulado - Sistema Runner",
    subcommands = {
        SignCommand.class,
        VerifyCommand.class
    }
)
public class Main implements Callable<Integer> {

    public static void main(String[] args) {
        int exitCode = new CommandLine(new Main()).execute(args);
        System.exit(exitCode);
    }

    @Override
    public Integer call() {
        // Caso o usuário execute apenas 'java -jar assinador.jar' sem subcomandos
        System.out.println("Use 'assinador --help' para ver os comandos disponíveis.");
        return 0;
    }
}