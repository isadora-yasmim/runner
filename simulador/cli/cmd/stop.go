package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o assinador.jar em modo servidor HTTP",
	Long: `Encerra o assinador.jar em execução na porta indicada.

Envia um sinal de parada via endpoint /stop do próprio servidor.

Exemplos:
  assinatura stop
  assinatura stop --port 9090`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isServerRunning(serverPort) {
			slog.Warn("nenhum assinador ativo na porta", "porta", serverPort)
			fmt.Fprintf(os.Stderr, "❌ Nenhum assinador em execução na porta %d.\n", serverPort)
			return nil
		}

		url := fmt.Sprintf("http://localhost:%d/stop", serverPort)
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			slog.Error("falha ao enviar sinal de parada", "url", url, "erro", err)
			return fmt.Errorf("erro ao enviar sinal de parada ao assinador: %w", err)
		}
		defer resp.Body.Close()

		slog.Info("assinador encerrado", "porta", serverPort)
		fmt.Printf("✅ Assinador encerrado na porta %d.\n", serverPort)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}