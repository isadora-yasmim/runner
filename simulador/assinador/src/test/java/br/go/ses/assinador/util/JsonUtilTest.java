// JsonUtilTest.java
package br.go.ses.assinador.util;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

import java.util.Map;

import static com.github.stefanbirkner.systemlambda.SystemLambda.*;

class JsonUtilTest {

    @Test
    void printSuccess_deveEscreverJsonNoStdout() throws Exception {
        String output = tapSystemOut(() ->
            JsonUtil.printSuccess("Operacao OK", Map.of("key", "value"))
        );
        assertTrue(output.contains("\"success\" : true"));
        assertTrue(output.contains("\"message\" : \"Operacao OK\""));
    }

    @Test
    void printError_deveEscreverJsonNoStderr() throws Exception {
        String output = tapSystemErr(() ->
            JsonUtil.printError("Algo deu errado")
        );
        assertTrue(output.contains("\"success\" : false"));
        assertTrue(output.contains("\"Algo deu errado\""));
    }
}