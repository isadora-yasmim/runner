package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type SignData struct {
	SignatureHash string `json:"signatureHash"`
	Algorithm     string `json:"algorithm"`
}

type ResponseOutput struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var signDocument string
var signTokenPin string
var signerJarPath string

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Simula a criação de uma assinatura digital",
	Run: func(cmd *cobra.Command, args []string) {

		cmdExec := exec.Command(
			"java",
			"-jar",
			signerJarPath,
			"sign",
			"-d", signDocument,
			"-t", signTokenPin,
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
			fmt.Println("✔ Assinatura criada com sucesso\n")

			dataBytes, _ := json.Marshal(resp.Data)
			var data SignData
			json.Unmarshal(dataBytes, &data)

			fmt.Println("→ Hash:", data.SignatureHash)
			fmt.Println("→ Algoritmo:", data.Algorithm)
		} else {
			fmt.Println("❌ Erro:")
			fmt.Println(resp.Message)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&signDocument, "document", "d", "", "Documento")
	signCmd.Flags().StringVarP(&signTokenPin, "token-pin", "t", "", "PIN")
	signCmd.Flags().StringVar(
		&signerJarPath,
		"jar",
		"..\\assinador\\target\\assinador-1.0-SNAPSHOT.jar",
		"Caminho do jar",
	)

	signCmd.MarkFlagRequired("document")
	signCmd.MarkFlagRequired("token-pin")
}