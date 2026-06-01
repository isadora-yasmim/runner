package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para operações de assinatura digital simulada",
	Long:  `assinatura — CLI do Sistema Runner para criação e validação de assinaturas digitais.`,
}

var serverPort int

// Execute é o ponto de entrada do CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(
		&serverPort,
		"port", "p",
		8080,
		"Porta do servidor HTTP do assinador",
	)
}