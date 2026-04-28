// ParameterValidatorTest.java
package br.go.ses.assinador.util;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class ParameterValidatorTest {

    // --- validateSign ---

    @Test
    void validateSign_documentoNulo_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateSign(null, "1234"));
    }

    @Test
    void validateSign_documentoVazio_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateSign("   ", "1234"));
    }

    @Test
    void validateSign_pinNulo_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateSign("doc.pdf", null));
    }

    @Test
    void validateSign_pinVazio_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateSign("doc.pdf", "  "));
    }

    @Test
    void validateSign_pinMenorQue4Chars_deveLancarExcecao() {
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateSign("doc.pdf", "123"));
        assertEquals("O PIN do token deve conter pelo menos 4 caracteres.", ex.getMessage());
    }

    @Test
    void validateSign_parametrosValidos_naoDeveLancarExcecao() {
        assertDoesNotThrow(() -> ParameterValidator.validateSign("doc.pdf", "1234"));
    }

    // --- validateVerify ---

    @Test
    void validateVerify_documentoNulo_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateVerify(null, "hash123"));
    }

    @Test
    void validateVerify_assinaturaNula_deveLancarExcecao() {
        assertThrows(IllegalArgumentException.class,
            () -> ParameterValidator.validateVerify("doc.pdf", null));
    }

    @Test
    void validateVerify_parametrosValidos_naoDeveLancarExcecao() {
        assertDoesNotThrow(() -> ParameterValidator.validateVerify("doc.pdf", "hash123"));
    }
}