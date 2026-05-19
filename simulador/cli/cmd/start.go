package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isServerRunning(serverPort) {
			fmt.Printf("✅ Assinador já está em execução na porta %d.\n", serverPort)
			return nil
		}

		jarPath := "..\\assinador\\target\\assinador-1.0-SNAPSHOT.jar"

		command := exec.Command(
			"java",
			"-jar",
			jarPath,
			"server",
			"--port",
			fmt.Sprintf("%d", serverPort),
		)

		err := command.Start()
		if err != nil {
			return fmt.Errorf("erro ao iniciar assinador.jar: %w", err)
		}

		fmt.Printf("✅ Assinador iniciado na porta %d.\n", serverPort)
		fmt.Printf("PID: %d\n", command.Process.Pid)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}