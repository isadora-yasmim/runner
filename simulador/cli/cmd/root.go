package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para operações de assinatura digital simulada",
	Long: `assinatura — CLI do Sistema Runner

Gerencia o ciclo de vida do assinador.jar e realiza operações de
criação e validação de assinaturas digitais simuladas.

Por padrão, as operações são realizadas via HTTP (modo servidor).
O assinador.jar é iniciado automaticamente quando necessário.

Exemplos:
  assinatura sign -d contrato.pdf -t 1234
  assinatura verify -d contrato.pdf -s <hash>
  assinatura start
  assinatura start --port 9090
  assinatura stop
  assinatura status
  assinatura version`,

	// PersistentPreRun executa após o parse de flags, antes de qualquer Run.
	// Garante que o logger esteja configurado com o nível correto.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

var (
	serverPort int
	verbose    bool
	quiet      bool
)

// Execute é o ponto de entrada do CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(
		&serverPort, "port", "p", 8080,
		"Porta do servidor HTTP do assinador",
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
		level = slog.LevelDebug // --verbose prevalece sobre --quiet
	case quiet:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}
