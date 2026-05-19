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

type SignRequest struct {
	Document string `json:"document"`
	TokenPin string `json:"tokenPin"`
}

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

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Simula a criação de uma assinatura digital",
	Run: func(cmd *cobra.Command, args []string) {

		err := ensureServerRunning()
		if err != nil {
			fmt.Println("❌ Erro ao garantir execução do assinador")
			fmt.Println(err)
			os.Exit(1)
		}		

		request := SignRequest{
			Document: signDocument,
			TokenPin: signTokenPin,
		}

		body, err := json.Marshal(request)
		if err != nil {
			fmt.Println("❌ Erro ao gerar requisição HTTP")
			os.Exit(1)
		}

		url := fmt.Sprintf("http://localhost:%d/sign", serverPort)

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
			fmt.Println("✔ Assinatura criada com sucesso\n")

			dataBytes, _ := json.Marshal(response.Data)

			var data SignData
			json.Unmarshal(dataBytes, &data)

			fmt.Println("→ Hash:", data.SignatureHash)
			fmt.Println("→ Algoritmo:", data.Algorithm)
		} else {
			fmt.Println("❌ Erro:")
			fmt.Println(response.Message)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(
		&signDocument,
		"document",
		"d",
		"",
		"Documento",
	)

	signCmd.Flags().StringVarP(
		&signTokenPin,
		"token-pin",
		"t",
		"",
		"PIN",
	)

	signCmd.MarkFlagRequired("document")
	signCmd.MarkFlagRequired("token-pin")
}