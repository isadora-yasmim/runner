package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var startPort int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isServerRunning(startPort) {
			fmt.Printf("✅ Assinador já está em execução na porta %d.\n", startPort)
			return nil
		}

		jarPath := "..\\assinador\\target\\assinador-1.0-SNAPSHOT.jar"

		command := exec.Command(
			"java",
			"-jar",
			jarPath,
			"server",
			"--port",
			fmt.Sprintf("%d", startPort),
		)

		err := command.Start()
		if err != nil {
			return fmt.Errorf("erro ao iniciar assinador.jar: %w", err)
		}

		fmt.Printf("✅ Assinador iniciado na porta %d.\n", startPort)
		fmt.Printf("PID: %d\n", command.Process.Pid)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&startPort, "port", "p", 8080, "Porta do servidor HTTP do assinador")
}