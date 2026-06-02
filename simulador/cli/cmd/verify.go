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
var verifyLocal bool

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Valida uma assinatura digital",
	Long: `Valida uma assinatura digital contra o documento informado.

Por padrão usa o modo HTTP (requer assinador em execução ou inicia automaticamente).
Use --local para invocar o assinador.jar diretamente via subprocess, sem servidor HTTP.

Exemplos:
  simulador verify -d contrato.pdf -s <hash>
  simulador verify -d contrato.pdf -s <hash> --local
  simulador verify -d contrato.pdf -s <hash> --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if verifyLocal {
			runVerifyLocal()
		} else {
			runVerifyHTTP()
		}
	},
}

func runVerifyLocal() {
	stdout, code, err := runLocalJar(
		"verify",
		"--document", verifyDocument,
		"--signature", verifySignature,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao executar assinador:", err)
		os.Exit(2)
	}

	fmt.Print(stdout)
	os.Exit(code)
}

func runVerifyHTTP() {
	if err := ensureServerRunning(); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Não foi possível conectar ao assinador HTTP.")
		fmt.Fprintln(os.Stderr, "→", err)
		os.Exit(2)
	}

	request := VerifyRequest{
		Document:  verifyDocument,
		Signature: verifySignature,
	}

	body, err := json.Marshal(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao preparar a requisição de validação.")
		os.Exit(2)
	}

	url := fmt.Sprintf("http://localhost:%d/validate", serverPort)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao enviar requisição para o assinador HTTP.")
		fmt.Fprintln(os.Stderr, "→", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao ler resposta do servidor.")
		os.Exit(2)
	}

	var response ResponseOutput
	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Erro ao interpretar resposta do assinador HTTP.")
		fmt.Fprintln(os.Stderr, string(responseBody))
		os.Exit(2)
	}

	if response.Success {
		fmt.Println("✔ Assinatura válida")
		return
	}

	fmt.Fprintln(os.Stderr, "❌ Assinatura inválida")
	fmt.Fprintln(os.Stderr, "→", response.Message)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(
		&verifyDocument,
		"document",
		"d",
		"",
		"Documento original",
	)

	verifyCmd.Flags().StringVarP(
		&verifySignature,
		"signature",
		"s",
		"",
		"Hash ou conteúdo da assinatura",
	)

	verifyCmd.Flags().BoolVar(
		&verifyLocal,
		"local",
		false,
		"Invoca o assinador.jar diretamente via subprocess (sem servidor HTTP)",
	)

	verifyCmd.MarkFlagRequired("document")
	verifyCmd.MarkFlagRequired("signature")
}