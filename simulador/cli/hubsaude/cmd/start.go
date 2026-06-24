package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/isadora-yasmim/simulador/internal/hubsaude"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o Simulador do HubSaúde",
	Long: `Inicia o simulador.jar em background na porta indicada.

Antes de iniciar, verifica se a porta está livre e garante que o JAR
esteja disponível em cache (~/.hubsaude/simulador), baixando-o do GitHub
Releases quando necessário.

Exemplos:
  hubsaude start
  hubsaude start --port 9595`,
	// Run não retorna error: controlamos o exit code manualmente para
	// distinguir erro de usuário (1) de erro de sistema (2).
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("resolvendo release mais recente do simulador")

		rel, err := hubsaude.LatestRelease()
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Não foi possível resolver a release do simulador.")
			fmt.Fprintln(os.Stderr, "→", err)
			os.Exit(2)
		}
		slog.Info("release resolvida", "versao", rel.Version, "jar", rel.JarName)

		jarPath, err := hubsaude.EnsureJar(rel)
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Falha ao obter o simulador.jar.")
			fmt.Fprintln(os.Stderr, "→", err)
			os.Exit(2)
		}

		pid, err := hubsaude.Start(jarPath, simPort)
		if err != nil {
			// Porta ocupada / já em execução = erro do usuário (exit 1).
			fmt.Fprintln(os.Stderr, "❌ Não foi possível iniciar o simulador.")
			fmt.Fprintln(os.Stderr, "→", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Simulador iniciado na porta %d (PID %d).\n", simPort, pid)
		fmt.Printf("→ Verifique a prontidão com: hubsaude status --port %d\n", simPort)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
