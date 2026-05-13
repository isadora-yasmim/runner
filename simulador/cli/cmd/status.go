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
			fmt.Printf("✅ Assinador está em execução na porta %d.\n", serverPort)
		} else {
			fmt.Printf("❌ Assinador não está em execução na porta %d.\n", serverPort)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}