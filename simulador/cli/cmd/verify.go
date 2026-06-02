package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// VerifyRequest representa o corpo da requisição POST /validate.
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
	Long: `Valida uma assinatura digital via assinador.jar.

Por padrão, usa o modo HTTP: o assinador.jar é iniciado automaticamente
se não estiver em execução.

Use --local para invocar o assinador.jar diretamente via subprocess,
sem iniciar o servidor HTTP.

Exit codes:
  0  Assinatura válida
  1  Assinatura inválida ou parâmetros rejeitados pelo assinador
  2  Erro de sistema (JVM ausente, JAR não encontrado, falha de rede)

Exemplos:
  assinatura verify -d contrato.pdf -s <hash>
  assinatura verify -d contrato.pdf -s <hash> --local
  assinatura verify -d contrato.pdf -s <hash> --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if verifyLocal {
			runVerifyLocal()
			return
		}

		runVerifyHTTP()
	},
}

func runVerifyLocal() {
	stdout, code, err := runLocalJar(
		"verify",
		"--document", verifyDocument,
		"--signature", verifySignature,
	)

	if err != nil {
		slog.Error("falha ao executar assinador em modo local", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível executar o assinador em modo local.\n→ %s\n", err)
		os.Exit(2)
	}

	fmt.Print(stdout)
	os.Exit(code)
}

func runVerifyHTTP() {
	if err := ensureServerRunning(); err != nil {
		slog.Error("falha ao iniciar o assinador", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível iniciar o assinador.\n→ %s\n", err)
		os.Exit(2)
	}

	request := VerifyRequest{Document: verifyDocument, Signature: verifySignature}
	body, err := json.Marshal(request)
	if err != nil {
		slog.Error("falha ao serializar requisição de verificação", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: falha interna ao preparar a requisição.\n→ %s\n", err)
		os.Exit(2)
	}

	url := fmt.Sprintf("http://localhost:%d/validate", serverPort)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("falha ao conectar ao assinador HTTP", "url", url, "erro", err)
		fmt.Fprintf(os.Stderr,
			"❌ Erro de sistema: não foi possível conectar ao assinador HTTP.\n→ %s\n→ Verifique com 'assinatura status' e reinicie com 'assinatura start'\n",
			err,
		)
		os.Exit(2)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("falha ao ler resposta do assinador", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: falha ao ler resposta do assinador.\n→ %s\n", err)
		os.Exit(2)
	}

	var response ResponseOutput
	if err := json.Unmarshal(responseBody, &response); err != nil {
		slog.Error("resposta inesperada do assinador", "corpo", string(responseBody), "erro", err)
		fmt.Fprintf(os.Stderr,
			"❌ Erro de sistema: o assinador retornou uma resposta inválida.\n→ Corpo recebido: %s\n",
			string(responseBody),
		)
		os.Exit(2)
	}

	if response.Success {
		fmt.Println("✔ Assinatura válida")
		return
	}

	slog.Info("assinatura inválida", "mensagem", response.Message)
	fmt.Fprintf(os.Stderr, "❌ Assinatura inválida: %s\n", response.Message)
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
		"Hash ou conteúdo da assinatura a validar",
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