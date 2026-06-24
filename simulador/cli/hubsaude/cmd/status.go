package cmd

import (
	"fmt"

	"github.com/isadora-yasmim/simulador/internal/hubsaude"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe o status do Simulador do HubSaúde",
	Long: `Exibe o estado do simulador distinguindo dois conceitos:

  • Liveness (saúde): o processo subiu e está vivo.
  • Readiness (prontidão): o serviço responde e está pronto para
    receber requisições (via /actuator/health, com fallback para
    porta TCP aberta).

Exemplos:
  hubsaude status
  hubsaude status --port 9595`,
	Run: func(cmd *cobra.Command, args []string) {
		st := hubsaude.Query()

		if !st.Registered {
			fmt.Println("❌ Nenhum simulador registrado como em execução.")
			fmt.Println("→ Use 'hubsaude start' para iniciar.")
			return
		}

		fmt.Printf("Simulador registrado: PID %d, porta %d\n", st.PID, st.Port)

		if st.Alive {
			fmt.Println("✅ Liveness: processo em execução (subiu).")
		} else {
			fmt.Println("❌ Liveness: processo não está mais em execução.")
			fmt.Println("→ O registro pode estar obsoleto; use 'hubsaude start'.")
			return
		}

		if st.Ready {
			fmt.Println("✅ Readiness: pronto para receber requisições.")
		} else {
			fmt.Println("⏳ Readiness: ainda não está pronto (inicializando).")
			fmt.Println("→ Aguarde alguns instantes e consulte novamente.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
