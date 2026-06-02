package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version e Commit são injetados via ldflags no build:
//
//	go build -ldflags "-X github.com/kyriosdata/assinatura/cmd.Version=v0.1.0 \
//	                   -X github.com/kyriosdata/assinatura/cmd.Commit=$(git rev-parse --short HEAD)"
//
// Em desenvolvimento (go run), exibem os valores padrão abaixo.
var (
	Version = "dev"
	Commit  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI",
	Long:  `Exibe a versão semântica e o SHA do commit do build atual.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("assinatura %s (%s)\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}