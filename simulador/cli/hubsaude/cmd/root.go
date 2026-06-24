package cmd

import (
	"log/slog"
	"os"

	"github.com/isadora-yasmim/simulador/internal/hubsaude"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hubsaude",
	Short: "CLI para gerenciar o ciclo de vida do Simulador do HubSaúde",
	Long: `hubsaude — CLI do Sistema Runner

Gerencia o Simulador do HubSaúde (simulador.jar): inicia, encerra e
consulta o status do processo. O JAR é obtido dinamicamente do GitHub
Releases e mantido em cache em ~/.hubsaude/simulador.

Exemplos:
  hubsaude start
  hubsaude start --port 9595
  hubsaude stop
  hubsaude status
  hubsaude version`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

var (
	simPort int
	verbose bool
	quiet   bool
)

// Execute é o ponto de entrada do binário hubsaude.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(
		&simPort, "port", "p", hubsaude.DefaultPort,
		"Porta do Simulador do HubSaúde",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Habilita saída detalhada (nível DEBUG)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&quiet, "quiet", false,
		"Suprime saída diagnóstica; exibe apenas erros",
	)
}

func initLogger() {
	var level slog.Level
	switch {
	case verbose:
		level = slog.LevelDebug
	case quiet:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}
