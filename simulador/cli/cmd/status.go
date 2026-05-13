package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusPort int

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica se o assinador.jar está em execução",
	Run: func(cmd *cobra.Command, args []string) {
		if isServerRunning(statusPort) {
			fmt.Printf("✅ Assinador está em execução na porta %d.\n", statusPort)
		} else {
			fmt.Printf("❌ Assinador não está em execução na porta %d.\n", statusPort)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().IntVarP(
		&statusPort,
		"port",
		"p",
		8080,
		"Porta do servidor HTTP do assinador",
	)
}