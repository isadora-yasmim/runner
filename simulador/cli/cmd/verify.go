package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var verifyDocument string
var verifySignature string
var verifyJarPath string

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Valida uma assinatura digital",
	Run: func(cmd *cobra.Command, args []string) {

		cmdExec := exec.Command(
			"java",
			"-jar",
			verifyJarPath,
			"verify",
			"-d", verifyDocument,
			"-s", verifySignature,
		)

		output, err := cmdExec.CombinedOutput()

		if err != nil {
			fmt.Println("❌ Erro ao executar assinador:")
			fmt.Println(string(output))
			os.Exit(1)
		}

		var resp ResponseOutput
		if err := json.Unmarshal(output, &resp); err != nil {
			fmt.Println("❌ Erro ao interpretar resposta")
			fmt.Println(string(output))
			os.Exit(1)
		}

		if resp.Success {
			fmt.Println("✔ Assinatura válida")
		} else {
			fmt.Println("❌ Assinatura inválida")
			fmt.Println("→", resp.Message)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&verifyDocument, "document", "d", "", "Documento")
	verifyCmd.Flags().StringVarP(&verifySignature, "signature", "s", "", "Assinatura")
	verifyCmd.Flags().StringVar(
		&verifyJarPath,
		"jar",
		"..\\assinador\\target\\assinador-1.0-SNAPSHOT.jar",
		"Caminho do jar",
	)

	verifyCmd.MarkFlagRequired("document")
	verifyCmd.MarkFlagRequired("signature")
}