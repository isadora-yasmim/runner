package br.go.ses.assinador.util;
import com.fasterxml.jackson.databind.ObjectMapper;
import br.go.ses.assinador.model.ResponseOutput;

public class JsonUtil {
    private static final ObjectMapper mapper = new ObjectMapper();

    public static void printSuccess(String message, Object data) {
        try {
            ResponseOutput response = new ResponseOutput(true, message, data);
            System.out.println(mapper.writerWithDefaultPrettyPrinter().writeValueAsString(response));
        } catch (Exception e) {
            System.err.println("Erro ao serializar resposta: " + e.getMessage());
        }
    }

    public static void printError(String message) {
        try {
            ResponseOutput response = new ResponseOutput(false, message, null);
            System.err.println(mapper.writerWithDefaultPrettyPrinter().writeValueAsString(response));
        } catch (Exception e) {
            System.err.println("Erro ao serializar erro: " + e.getMessage());
        }
    }
}