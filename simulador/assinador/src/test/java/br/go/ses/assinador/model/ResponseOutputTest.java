// ResponseOutputTest.java

package br.go.ses.assinador.model;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class ResponseOutputTest {

    @Test
    void construtor_devePreencherTodosCampos() {
        ResponseOutput out = new ResponseOutput(true, "ok", "data");
        assertTrue(out.isSuccess());
        assertEquals("ok", out.getMessage());
        assertEquals("data", out.getData());
    }

    @Test
    void construtor_comDadoNulo_deveFuncionar() {
        ResponseOutput out = new ResponseOutput(false, "erro", null);
        assertFalse(out.isSuccess());
        assertNull(out.getData());
    }
}