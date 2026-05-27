package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica se o assinador.jar está em execução",
	Run: func(cmd *cobra.Command, args []string) {
		if isServerRunning(serverPort) {
			fmt.Printf("✅ Assinador HTTP ativo na porta %d.\n", serverPort)
			fmt.Println("→ O CLI reutilizará esta instância nas próximas operações.")
		} else {
			fmt.Printf("❌ Nenhum assinador HTTP ativo na porta %d.\n", serverPort)
			fmt.Println("→ Os comandos sign e verify tentarão iniciar o assinador automaticamente.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}