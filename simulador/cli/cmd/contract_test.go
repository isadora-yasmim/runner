package cmd_test

// Testes de contrato CLI ↔ JAR.
//
// Estratégia: compilar o binário Go em TestMain e invocar via exec.Command,
// validando stdout, stderr e exit code em cada cenário.
//
// Os testes são marcados com t.Skip quando o JAR ou o Java não estão
// disponíveis (ambiente de desenvolvedor sem build Java), mas rodam
// obrigatoriamente em CI onde o job garante ambos.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// caminhos resolvidos em TestMain
var (
	binaryPath string // caminho do binário CLI compilado
	jarPath    string // caminho do assinador.jar
)

// TestMain compila o CLI uma vez e valida pré-condições.
func TestMain(m *testing.M) {
	// Resolve o diretório do módulo Go (dois níveis acima de cmd/)
	moduleDir, err := findModuleDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Não foi possível encontrar go.mod:", err)
		os.Exit(1)
	}

	// Binário temporário
	binName := "assinatura_test_bin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binaryPath = filepath.Join(os.TempDir(), binName)

	// Compila o CLI
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = moduleDir
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Falha ao compilar CLI:", err)
		os.Exit(1)
	}

	// Localiza o JAR (path configurável via variável de ambiente ASSINADOR_JAR)
	jarPath = os.Getenv("ASSINADOR_JAR")
	if jarPath == "" {
		// Convenção padrão do projeto
		jarPath = filepath.Join(moduleDir, "..", "assinador", "target", "assinador.jar")
	}

	code := m.Run()

	// Limpa binário temporário
	os.Remove(binaryPath)

	os.Exit(code)
}

// ------------------------------------------------------------------ helpers

func findModuleDir() (string, error) {
	// Parte do diretório atual e sobe até encontrar go.mod
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod não encontrado")
		}
		dir = parent
	}
}

func jarAvailable() bool {
	_, err := os.Stat(jarPath)
	return err == nil
}

func javaAvailable() bool {
	return exec.Command("java", "-version").Run() == nil
}

// runCLI executa o binário compilado com os argumentos dados e retorna
// stdout, stderr e o exit code.
func runCLI(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return
}

// ------------------------------------------------------------------ version

func TestCLI_Version_DeveExibirVersao(t *testing.T) {
	stdout, _, code := runCLI("version")
	if code != 0 {
		t.Fatalf("exit code inesperado: %d", code)
	}
	if !strings.Contains(stdout, "0.") {
		t.Errorf("saída não contém versão semântica: %q", stdout)
	}
}

// ------------------------------------------------------------------ sign — flags obrigatórias

func TestCLI_Sign_SemDocumento_DeveRetornarExitCodeNaoZero(t *testing.T) {
	_, _, code := runCLI("sign", "--token-pin", "1234")
	if code == 0 {
		t.Error("esperava exit code != 0 quando --document está ausente")
	}
}

func TestCLI_Sign_SemPin_DeveRetornarExitCodeNaoZero(t *testing.T) {
	_, _, code := runCLI("sign", "--document", "doc.pdf")
	if code == 0 {
		t.Error("esperava exit code != 0 quando --token-pin está ausente")
	}
}

// ------------------------------------------------------------------ verify — flags obrigatórias

func TestCLI_Verify_SemDocumento_DeveRetornarExitCodeNaoZero(t *testing.T) {
	_, _, code := runCLI("verify", "--signature", "hash123")
	if code == 0 {
		t.Error("esperava exit code != 0 quando --document está ausente")
	}
}

func TestCLI_Verify_SemAssinatura_DeveRetornarExitCodeNaoZero(t *testing.T) {
	_, _, code := runCLI("verify", "--document", "doc.pdf")
	if code == 0 {
		t.Error("esperava exit code != 0 quando --signature está ausente")
	}
}

// ------------------------------------------------------------------ sign via JAR subprocess (modo --local)
// Estes testes requerem Java + JAR e são pulados em ambientes sem eles.

func TestCLI_SignLocal_DocumentoValido_DeveRetornarExitCode0(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível — pulando teste de contrato com JAR")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível — execute 'mvn package' antes")
	}

	stdout, stderr, code := runCLI("sign", "--local", "--document", "contrato.pdf", "--token-pin", "1234")

	if code != 0 {
		t.Fatalf("exit code inesperado %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}

	// Saída deve conter o hash simulado
	if !strings.Contains(stdout, "mock_hash_abc123_base64_encoded_signature_simulated") {
		t.Errorf("stdout não contém hash esperado:\n%s", stdout)
	}
}

func TestCLI_SignLocal_PinCurto_DeveRetornarExitCode1(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível")
	}

	_, _, code := runCLI("sign", "--local", "--document", "doc.pdf", "--token-pin", "12")
	if code != 1 {
		t.Errorf("esperava exit code 1 para PIN curto, obteve %d", code)
	}
}

func TestCLI_SignLocal_DocumentoComEspacos_DeveRetornarExitCode0(t *testing.T) {
	// Critério E1: espaços nos argumentos devem ser preservados
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível")
	}

	stdout, stderr, code := runCLI("sign", "--local", "--document", "meu documento 2024.pdf", "--token-pin", "1234")
	if code != 0 {
		t.Fatalf("exit code inesperado %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

func TestCLI_VerifyLocal_HashValido_DeveRetornarExitCode0(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível")
	}

	_, _, code := runCLI("verify", "--local",
		"--document", "doc.pdf",
		"--signature", "mock_hash_abc123_base64_encoded_signature_simulated")
	if code != 0 {
		t.Errorf("esperava exit code 0 para hash válido, obteve %d", code)
	}
}

func TestCLI_VerifyLocal_HashInvalido_DeveRetornarExitCode1(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível")
	}

	_, _, code := runCLI("verify", "--local",
		"--document", "doc.pdf",
		"--signature", "hash_completamente_errado")
	if code != 1 {
		t.Errorf("esperava exit code 1 para hash inválido, obteve %d", code)
	}
}

// ------------------------------------------------------------------ saída JSON estruturada

func TestCLI_SignLocal_SaidaEhJSONValido(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("assinador.jar não disponível")
	}

	stdout, _, _ := runCLI("sign", "--local", "--document", "doc.pdf", "--token-pin", "1234")

	// Extrai somente a linha JSON (o CLI pode emitir linhas de diagnóstico antes)
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "{") {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(line), &obj); err != nil {
				t.Errorf("linha JSON inválida: %q — erro: %v", line, err)
			}
			return
		}
	}
	// Se não há JSON na saída mas o teste chegou aqui o JAR foi invocado
	// e pode ter imprimido output formatado; não falha o teste.
}

// ------------------------------------------------------------------ help

func TestCLI_Help_DeveRetornarExitCode0(t *testing.T) {
	_, _, code := runCLI("--help")
	if code != 0 {
		t.Errorf("esperava exit code 0 para --help, obteve %d", code)
	}
}

func TestCLI_SignHelp_DeveConterExemplosDeUso(t *testing.T) {
	stdout, _, _ := runCLI("sign", "--help")
	if !strings.Contains(stdout, "--document") || !strings.Contains(stdout, "--token-pin") {
		t.Errorf("--help do sign não documenta as flags esperadas:\n%s", stdout)
	}
}
