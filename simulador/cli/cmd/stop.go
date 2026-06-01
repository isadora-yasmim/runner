package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o assinador.jar em modo servidor HTTP",
	Long: `Encerra o assinador.jar em execução na porta indicada.

Envia um sinal de parada via endpoint /stop e valida a resposta recebida
antes de confirmar o encerramento.

Exit codes:
  0  Servidor encerrado com sucesso
  1  Servidor reportou falha ao encerrar
  2  Erro de sistema (falha de rede, resposta inválida)

Exemplos:
  assinatura stop
  assinatura stop --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isServerRunning(serverPort) {
			slog.Warn("nenhum assinador ativo na porta", "porta", serverPort)
			fmt.Fprintf(os.Stderr,
				"❌ Nenhum assinador em execução na porta %d.\n→ Use 'assinatura status' para verificar.\n",
				serverPort,
			)
			return
		}

		url := fmt.Sprintf("http://localhost:%d/stop", serverPort)
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			slog.Error("falha ao enviar sinal de parada", "url", url, "erro", err)
			fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível enviar sinal de parada ao assinador.\n→ %s\n", err)
			os.Exit(2)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("falha ao ler resposta do /stop", "erro", err)
			fmt.Fprintf(os.Stderr, "❌ Erro de sistema: falha ao ler resposta do endpoint /stop.\n→ %s\n", err)
			os.Exit(2)
		}

		var response ResponseOutput
		if err := json.Unmarshal(body, &response); err != nil {
			slog.Error("resposta inesperada do /stop", "corpo", string(body), "erro", err)
			fmt.Fprintf(os.Stderr,
				"❌ Erro de sistema: resposta inesperada do servidor (não é JSON válido).\n→ Corpo: %s\n",
				string(body),
			)
			os.Exit(2)
		}

		if !response.Success {
			slog.Error("servidor reportou falha ao encerrar", "mensagem", response.Message)
			fmt.Fprintf(os.Stderr, "❌ O servidor reportou falha ao encerrar: %s\n", response.Message)
			os.Exit(1)
		}

		slog.Info("assinador encerrado com sucesso", "porta", serverPort)
		fmt.Printf("✅ Assinador encerrado na porta %d.\n", serverPort)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}