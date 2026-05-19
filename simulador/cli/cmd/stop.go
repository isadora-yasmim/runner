package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o assinador.jar em modo servidor HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isServerRunning(serverPort) {
			fmt.Printf("❌ Nenhum assinador em execução na porta %d.\n", serverPort)
			return nil
		}

		url := fmt.Sprintf("http://localhost:%d/stop", serverPort)

		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			return fmt.Errorf("erro ao encerrar assinador: %w", err)
		}
		defer resp.Body.Close()

		fmt.Printf("✅ Assinador encerrado na porta %d.\n", serverPort)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}