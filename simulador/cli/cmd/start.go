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

		err := checkJavaInstalled()
		if err != nil {
			return err
		}

		jarPath, err := getAssinadorJarPath()
		if err != nil {
			return err
		}

		command := exec.Command(
			"java",
			"-jar",
			jarPath,
			"server",
			"--port",
			fmt.Sprintf("%d", serverPort),
		)

		err = command.Start()
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