package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type VerifyRequest struct {
	Document  string `json:"document"`
	Signature string `json:"signature"`
}

var verifyDocument string
var verifySignature string

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Valida uma assinatura digital",
	Run: func(cmd *cobra.Command, args []string) {

		err := ensureServerRunning()
		if err != nil {
			fmt.Println("❌ Erro ao garantir execução do assinador")
			fmt.Println(err)
			os.Exit(1)
		}

		request := VerifyRequest{
			Document:  verifyDocument,
			Signature: verifySignature,
		}

		body, err := json.Marshal(request)
		if err != nil {
			fmt.Println("❌ Erro ao gerar requisição HTTP")
			os.Exit(1)
		}

		url := fmt.Sprintf(
			"http://localhost:%d/validate",
			serverPort,
		)

		resp, err := http.Post(
			url,
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			fmt.Println("❌ Erro ao conectar com o assinador HTTP")
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("❌ Erro ao ler resposta do servidor")
			os.Exit(1)
		}

		var response ResponseOutput

		if err := json.Unmarshal(responseBody, &response); err != nil {
			fmt.Println("❌ Erro ao interpretar resposta")
			fmt.Println(string(responseBody))
			os.Exit(1)
		}

		if response.Success {
			fmt.Println("✔ Assinatura válida")
		} else {
			fmt.Println("❌ Assinatura inválida")
			fmt.Println("→", response.Message)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(
		&verifyDocument,
		"document",
		"d",
		"",
		"Documento",
	)

	verifyCmd.Flags().StringVarP(
		&verifySignature,
		"signature",
		"s",
		"",
		"Assinatura",
	)

	verifyCmd.MarkFlagRequired("document")
	verifyCmd.MarkFlagRequired("signature")
}