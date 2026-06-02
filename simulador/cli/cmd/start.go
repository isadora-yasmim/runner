package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	Long: `Inicia o assinador.jar como servidor HTTP em background.

Realiza um health check real antes de iniciar: se uma instância já
estiver ativa na porta indicada, o comando encerra sem criar duplicatas.

Exit codes:
  0  Servidor iniciado com sucesso (ou já estava em execução)
  2  Erro de sistema (JVM ausente, JAR não encontrado, falha ao iniciar processo)

Exemplos:
  assinatura start
  assinatura start --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if isServerRunning(serverPort) {
			slog.Info("assinador já em execução", "porta", serverPort)
			fmt.Printf("✅ Assinador já está em execução na porta %d.\n", serverPort)
			return
		}

		if err := checkJavaInstalled(); err != nil {
			slog.Error("Java não disponível", "erro", err)
			fmt.Fprintf(os.Stderr, "❌ Erro de sistema: %s\n", err)
			os.Exit(2)
		}

		jarPath, err := getAssinadorJarPath()
		if err != nil {
			slog.Error("JAR não encontrado", "erro", err)
			fmt.Fprintf(os.Stderr, "❌ Erro de sistema: %s\n", err)
			os.Exit(2)
		}

		command := exec.Command("java", "-jar", jarPath, "server", "--port", fmt.Sprintf("%d", serverPort))
		if err := command.Start(); err != nil {
			slog.Error("falha ao iniciar assinador.jar", "jar", jarPath, "erro", err)
			fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível iniciar o processo do assinador.\n→ %s\n", err)
			os.Exit(2)
		}

		slog.Info("assinador iniciado", "pid", command.Process.Pid, "porta", serverPort)
		fmt.Printf("✅ Assinador iniciado na porta %d (PID %d).\n", serverPort, command.Process.Pid)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}