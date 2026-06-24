package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version e Commit são injetados em tempo de build via -ldflags:
//
//	-X '.../hubsaude/cmd.Version=v1.0.0' -X '.../hubsaude/cmd.Commit=abc1234'
//
// Os valores padrão indicam build local (sem release).
var (
	Version = "dev"
	Commit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI hubsaude",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hubsaude %s (%s)\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
