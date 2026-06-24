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

// SignRequest representa o corpo da requisição POST /sign.
type SignRequest struct {
	Document string `json:"document"`
	TokenPin string `json:"tokenPin"`
}

// SignData contém os campos retornados em caso de assinatura bem-sucedida.
type SignData struct {
	SignatureHash string `json:"signatureHash"`
	Algorithm     string `json:"algorithm"`
}

var signDocument string
var signTokenPin string
var signLocal bool

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Simula a criação de uma assinatura digital",
	Long: `Cria uma assinatura digital simulada via assinador.jar.

Por padrão, usa o modo HTTP: o assinador.jar é iniciado automaticamente
se não estiver em execução.

Use --local para invocar o assinador.jar diretamente via subprocess,
sem iniciar o servidor HTTP.

Exit codes:
  0  Assinatura criada com sucesso
  1  Erro do usuário (parâmetros rejeitados pelo assinador)
  2  Erro de sistema (JVM ausente, JAR não encontrado, falha de rede)

Exemplos:
  assinatura sign -d contrato.pdf -t 1234
  assinatura sign -d contrato.pdf -t 1234 --local
  assinatura sign -d contrato.pdf -t 1234 --port 9090`,
	Run: func(cmd *cobra.Command, args []string) {
		if signLocal {
			runSignLocal()
			return
		}

		runSignHTTP()
	},
}

func runSignLocal() {
	stdout, code, err := runLocalJar(
		"sign",
		"--document", signDocument,
		"--token-pin", signTokenPin,
	)

	if err != nil {
		slog.Error("falha ao executar assinador em modo local", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível executar o assinador em modo local.\n→ %s\n", err)
		os.Exit(2)
	}

	fmt.Print(stdout)
	os.Exit(code)
}

func runSignHTTP() {
	if err := ensureServerRunning(); err != nil {
		slog.Error("falha ao iniciar o assinador", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: não foi possível iniciar o assinador.\n→ %s\n", err)
		os.Exit(2)
	}

	request := SignRequest{Document: signDocument, TokenPin: signTokenPin}
	body, err := json.Marshal(request)
	if err != nil {
		slog.Error("falha ao serializar requisição", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: falha interna ao preparar a requisição.\n→ %s\n", err)
		os.Exit(2)
	}

	url := fmt.Sprintf("http://localhost:%d/sign", serverPort)
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

	if !response.Success {
		slog.Info("assinatura recusada pelo assinador", "mensagem", response.Message)
		fmt.Fprintf(os.Stderr, "❌ Parâmetros inválidos: %s\n", response.Message)
		os.Exit(1)
	}

	dataBytes, _ := json.Marshal(response.Data)
	var data SignData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		slog.Error("falha ao interpretar dados da assinatura", "erro", err)
		fmt.Fprintf(os.Stderr, "❌ Erro de sistema: falha ao interpretar dados da assinatura retornada.\n→ %s\n", err)
		os.Exit(2)
	}

	fmt.Println("✔ Assinatura criada com sucesso")
	fmt.Println()
	fmt.Printf("→ Hash:      %s\n", data.SignatureHash)
	fmt.Printf("→ Algoritmo: %s\n", data.Algorithm)
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(
		&signDocument,
		"document",
		"d",
		"",
		"Documento a ser assinado (caminho ou conteúdo)",
	)

	signCmd.Flags().StringVarP(
		&signTokenPin,
		"token-pin",
		"t",
		"",
		"PIN do dispositivo criptográfico (mínimo 4 caracteres)",
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
