package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var startTimeout int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	Long: `Inicia o assinador.jar em modo servidor HTTP em background.

O servidor responde nos endpoints /sign, /validate, /health e /stop.

Exemplos:
  assinatura start                 → inicia na porta padrao (8080)
  assinatura start --port 8081     → inicia na porta 8081
  assinatura start --timeout 10    → encerra apos 10 min de inatividade`,
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

		jarArgs := []string{
			"-jar",
			jarPath,
			"server",
			"--port",
			fmt.Sprintf("%d", serverPort),
		}

		if startTimeout > 0 {
			jarArgs = append(jarArgs, "--timeout", fmt.Sprintf("%d", startTimeout))
		}

		command := exec.Command("java", jarArgs...)

		err = command.Start()
		if err != nil {
			return fmt.Errorf("erro ao iniciar assinador.jar: %w", err)
		}

		fmt.Printf("✅ Assinador iniciado na porta %d.\n", serverPort)
		fmt.Printf("PID: %d\n", command.Process.Pid)
		if startTimeout > 0 {
			fmt.Printf("⏱️  Auto-shutdown após %d min de inatividade.\n", startTimeout)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVar(
		&startTimeout,
		"timeout",
		0,
		"Minutos de inatividade antes do auto-shutdown (0 = desativado)",
	)
}