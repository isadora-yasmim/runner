package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica se o assinador.jar está em execução",
	Long: `Realiza um health check HTTP no assinador.jar e exibe o status atual.

Exemplos:
  assinatura status
  assinatura status --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if isServerRunning(serverPort) {
			slog.Debug("health check OK", "porta", serverPort)
			fmt.Printf("✅ Assinador HTTP ativo na porta %d.\n", serverPort)
			fmt.Println("→ Os comandos sign e verify usarão esta instância.")
		} else {
			slog.Debug("health check sem resposta", "porta", serverPort)
			fmt.Printf("❌ Nenhum assinador HTTP ativo na porta %d.\n", serverPort)
			fmt.Println("→ Use 'assinatura start' para iniciar o servidor.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
