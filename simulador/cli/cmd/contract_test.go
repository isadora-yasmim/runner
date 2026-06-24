package cmd_test

// Testes de contrato CLI ↔ JAR.
//
// Estratégia: compilar o binário Go em TestMain e invocar via exec.Command,
// validando stdout, stderr e exit code em cada cenário.
//
// Os testes que dependem do JAR fazem t.Skip quando ASSINADOR_JAR não está
// definido (ambiente sem build Java). Em CI o job 'contrato' define a variável
// apontando para o JAR compilado pelo job 'java'.

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

var binaryPath string // caminho do binário CLI compilado pelo TestMain

// TestMain compila o CLI uma vez e disponibiliza o binário para todos os testes.
func TestMain(m *testing.M) {
	moduleDir, err := findModuleDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Não foi possível encontrar go.mod:", err)
		os.Exit(1)
	}

	binName := "assinatura_test_bin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binaryPath = filepath.Join(os.TempDir(), binName)

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = moduleDir
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Falha ao compilar CLI:", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(binaryPath)
	os.Exit(code)
}

// ------------------------------------------------------------------ helpers

func findModuleDir() (string, error) {
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

// jarAvailable retorna true quando ASSINADOR_JAR está definido e o arquivo existe.
func jarAvailable() bool {
	env := os.Getenv("ASSINADOR_JAR")
	if env == "" {
		return false
	}
	_, err := os.Stat(env)
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

// extractJSON localiza e retorna o primeiro bloco JSON completo no stdout.
// Necessário porque o JAR usa pretty-print (Jackson writerWithDefaultPrettyPrinter),
// então o JSON ocupa múltiplas linhas e não pode ser parseado linha a linha.
func extractJSON(stdout string) (string, bool) {
	start := strings.Index(stdout, "{")
	if start == -1 {
		return "", false
	}
	// Encontra o fechamento correspondente contando profundidade
	depth := 0
	for i := start; i < len(stdout); i++ {
		switch stdout[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return stdout[start : i+1], true
			}
		}
	}
	return "", false
}

// ------------------------------------------------------------------ version

func TestCLI_Version_DeveExibirVersao(t *testing.T) {
	stdout, _, code := runCLI("version")
	if code != 0 {
		t.Fatalf("exit code inesperado: %d", code)
	}

	if !strings.Contains(stdout, "assinatura") {
		t.Errorf("saída não contém nome do CLI: %q", stdout)
	}

	if !strings.Contains(stdout, "(") || !strings.Contains(stdout, ")") {
		t.Errorf("saída não contém commit entre parênteses: %q", stdout)
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

// ------------------------------------------------------------------ sign --local (requerem JAR)

func TestCLI_SignLocal_DocumentoValido_DeveRetornarExitCode0(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
	}

	stdout, stderr, code := runCLI("sign", "--local", "--document", "contrato.pdf", "--token-pin", "1234")
	if code != 0 {
		t.Fatalf("exit code inesperado %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "mock_hash_abc123_base64_encoded_signature_simulated") {
		t.Errorf("stdout não contém hash esperado:\n%s", stdout)
	}
}

func TestCLI_SignLocal_PinCurto_DeveRetornarExitCode1(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
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
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
	}

	stdout, stderr, code := runCLI("sign", "--local", "--document", "meu documento 2024.pdf", "--token-pin", "1234")
	if code != 0 {
		t.Fatalf("exit code inesperado %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

// ------------------------------------------------------------------ verify --local (requerem JAR)

func TestCLI_VerifyLocal_HashValido_DeveRetornarExitCode0(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
	}

	stdout, stderr, code := runCLI(
		"verify", "--local",
		"--document", "doc.pdf",
		"--signature", "mock_hash_abc123_base64_encoded_signature_simulated",
	)
	if code != 0 {
		t.Fatalf("exit code inesperado %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

func TestCLI_VerifyLocal_HashInvalido_DeveRetornarExitCode1(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
	}

	_, _, code := runCLI(
		"verify", "--local",
		"--document", "doc.pdf",
		"--signature", "hash_completamente_errado",
	)
	if code != 1 {
		t.Errorf("esperava exit code 1 para hash inválido, obteve %d", code)
	}
}

// ------------------------------------------------------------------ saída JSON

func TestCLI_SignLocal_SaidaEhJSONValido(t *testing.T) {
	if !javaAvailable() {
		t.Skip("Java não disponível")
	}
	if !jarAvailable() {
		t.Skip("ASSINADOR_JAR não definido ou arquivo ausente")
	}

	stdout, _, code := runCLI("sign", "--local", "--document", "doc.pdf", "--token-pin", "1234")
	if code != 0 {
		t.Skipf("sign --local retornou exit code %d; pulando validação de JSON", code)
	}

	// O JAR usa pretty-print (múltiplas linhas): extrai o bloco JSON completo
	// em vez de tentar parsear linha por linha.
	jsonBlock, found := extractJSON(stdout)
	if !found {
		t.Fatalf("nenhum bloco JSON encontrado no stdout:\n%s", stdout)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonBlock), &obj); err != nil {
		t.Errorf("JSON inválido: %v\nbloco extraído: %s", err, jsonBlock)
	}
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
