package cmd

import (
	"fmt"
	"os"

	"github.com/isadora-yasmim/simulador/internal/hubsaude"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o Simulador do HubSaúde",
	Long: `Encerra o processo do simulador registrado em ~/.hubsaude/simulador.pid.

Exemplos:
  hubsaude stop`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := hubsaude.Stop()
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Não foi possível encerrar o simulador.")
			fmt.Fprintln(os.Stderr, "→", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Simulador encerrado (PID %d).\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
