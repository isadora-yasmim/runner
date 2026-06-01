package cmd

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	Long: `Inicia o assinador.jar como servidor HTTP em background.

Realiza um health check real antes de iniciar: se uma instância já
estiver ativa na porta indicada, o comando encerra sem criar duplicatas.

Exemplos:
  assinatura start
  assinatura start --port 9090`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isServerRunning(serverPort) {
			slog.Info("assinador já em execução", "porta", serverPort)
			fmt.Printf("✅ Assinador já está em execução na porta %d.\n", serverPort)
			return nil
		}

		if err := checkJavaInstalled(); err != nil {
			return err
		}

		jarPath, err := getAssinadorJarPath()
		if err != nil {
			return err
		}

		command := exec.Command("java", "-jar", jarPath, "server", "--port", fmt.Sprintf("%d", serverPort))
		if err := command.Start(); err != nil {
			return fmt.Errorf("erro ao iniciar assinador.jar: %w", err)
		}

		slog.Info("assinador iniciado", "pid", command.Process.Pid, "porta", serverPort)
		fmt.Printf("✅ Assinador iniciado na porta %d (PID %d).\n", serverPort, command.Process.Pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}