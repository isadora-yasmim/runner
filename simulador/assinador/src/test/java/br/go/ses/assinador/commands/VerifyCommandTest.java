// VerifyCommandTest.java

package br.go.ses.assinador.commands;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import picocli.CommandLine;

class VerifyCommandTest {

    private static final String HASH_VALIDO =
        "mock_hash_abc123_base64_encoded_signature_simulated";

    @Test
    void verify_hashValido_deveRetornarExitCode0() {
        int exitCode = new CommandLine(new VerifyCommand())
            .execute("-d", "doc.pdf", "-s", HASH_VALIDO);
        assertEquals(0, exitCode);
    }

    @Test
    void verify_hashInvalido_deveRetornarExitCode1() {
        int exitCode = new CommandLine(new VerifyCommand())
            .execute("-d", "doc.pdf", "-s", "hash_errado");
        assertEquals(1, exitCode);
    }

    @Test
    void verify_semAssinatura_naoDeveRetornarExitCode0() {
        int exitCode = new CommandLine(new VerifyCommand())
            .execute("-d", "doc.pdf");
        assertNotEquals(0, exitCode);
    }
}