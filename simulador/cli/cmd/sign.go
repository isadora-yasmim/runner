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
var signLocal bool

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Simula a criação de uma assinatura digital",
	Long: `Cria uma assinatura digital simulada para o documento informado.

Por padrão usa o modo HTTP (requer assinador em execução ou inicia automaticamente).
Use --local para invocar o assinador.jar diretamente via subprocess, sem servidor HTTP.

Exemplos:
  simulador sign -d contrato.pdf -t 1234
  simulador sign -d contrato.pdf -t 1234 --local
  simulador sign -d contrato.pdf -t 1234 --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if signLocal {
			runSignLocal()
		} else {
			runSignHTTP()
		}
	},
}

func runSignLocal() {
	stdout, code, err := runLocalJar(
		"sign",
		"--document", signDocument,
		"--token-pin", signTokenPin,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao executar assinador:", err)
		os.Exit(2)
	}

	fmt.Print(stdout)
	os.Exit(code)
}

func runSignHTTP() {
	if err := ensureServerRunning(); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao garantir execução do assinador")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	request := SignRequest{
		Document: signDocument,
		TokenPin: signTokenPin,
	}

	body, err := json.Marshal(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao gerar requisição HTTP")
		os.Exit(2)
	}

	url := fmt.Sprintf("http://localhost:%d/sign", serverPort)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao conectar com o assinador HTTP")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao ler resposta do servidor")
		os.Exit(2)
	}

	var response ResponseOutput
	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao interpretar resposta")
		fmt.Fprintln(os.Stderr, string(responseBody))
		os.Exit(2)
	}

	if response.Success {
		fmt.Println("✔ Assinatura criada com sucesso\n")

		dataBytes, _ := json.Marshal(response.Data)
		var data SignData
		json.Unmarshal(dataBytes, &data)

		fmt.Println("→ Hash:", data.SignatureHash)
		fmt.Println("→ Algoritmo:", data.Algorithm)
	} else {
		fmt.Fprintln(os.Stderr, "❌ Erro:")
		fmt.Fprintln(os.Stderr, response.Message)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(
		&signDocument,
		"document",
		"d",
		"",
		"Caminho ou conteúdo do documento a ser assinado",
	)

	signCmd.Flags().StringVarP(
		&signTokenPin,
		"token-pin",
		"t",
		"",
		"PIN do dispositivo criptográfico",
	)

	signCmd.Flags().BoolVar(
		&signLocal,
		"local",
		false,
		"Invoca o assinador.jar diretamente via subprocess (sem servidor HTTP)",
	)

	signCmd.MarkFlagRequired("document")
	signCmd.MarkFlagRequired("token-pin")
}