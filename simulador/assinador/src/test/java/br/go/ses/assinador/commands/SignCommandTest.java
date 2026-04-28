// SignCommandTest.java
package br.go.ses.assinador.commands;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import picocli.CommandLine;
import static com.github.stefanbirkner.systemlambda.SystemLambda.*;

class SignCommandTest {

    @Test
    void sign_parametrosValidos_deveRetornarExitCode0() throws Exception {
        int exitCode = new CommandLine(new SignCommand())
            .execute("-d", "documento.pdf", "-t", "1234");
        assertEquals(0, exitCode);
    }

    @Test
    void sign_semDocumento_deveRetornarExitCode1() {
        // picocli retorna 2 para parâmetros obrigatórios ausentes
        int exitCode = new CommandLine(new SignCommand())
            .execute("-t", "1234");
        assertNotEquals(0, exitCode);
    }

    @Test
    void sign_pinMuitoCurto_deveRetornarExitCode1() throws Exception {
        int exitCode = new CommandLine(new SignCommand())
            .execute("-d", "doc.pdf", "-t", "12");
        assertEquals(1, exitCode);
    }

    @Test
    void sign_parametrosValidos_deveProduzirJsonComHash() throws Exception {
        String output = tapSystemOut(() ->
            new CommandLine(new SignCommand()).execute("-d", "doc.pdf", "-t", "1234")
        );
        assertTrue(output.contains("mock_hash_abc123_base64_encoded_signature_simulated"));
        assertTrue(output.contains("SHA256withRSA"));
    }
}