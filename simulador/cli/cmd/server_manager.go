package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func ensureServerRunning() error {

	if isServerRunning(serverPort) {
		return nil
	}

	err := checkJavaInstalled()
	if err != nil {
		return err
	}

	jarPath, err := getAssinadorJarPath()
	if err != nil {
		return err
	}

	fmt.Printf("⚠️ Assinador não encontrado na porta %d.\n", serverPort)
	fmt.Println("🚀 Iniciando automaticamente...")

	command := exec.Command(
		"java",
		"-jar",
		jarPath,
		"server",
		"--port",
		fmt.Sprintf("%d", serverPort),
	)

	err = command.Start()
	if err != nil {
		return fmt.Errorf("erro ao iniciar assinador automaticamente: %w", err)
	}

	fmt.Println("✅ Assinador iniciado automaticamente.")
	fmt.Printf("PID: %d\n", command.Process.Pid)

	return waitForServer()
}

// runLocalJar invoca o assinador.jar diretamente via subprocess (modo --local).
// stdout do JAR é capturado e retornado; stderr é propagado para os.Stderr.
// O exit code do JAR é propagado via error retornado.
func runLocalJar(args ...string) (string, int, error) {
	jarPath, err := getAssinadorJarPath()
	if err != nil {
		return "", 2, err
	}

	if err := checkJavaInstalled(); err != nil {
		return "", 2, err
	}

	cmdArgs := append([]string{"-jar", jarPath}, args...)
	cmd := exec.Command("java", cmdArgs...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	stdout := string(out)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdout, exitErr.ExitCode(), nil
		}
		return stdout, 2, fmt.Errorf("erro ao executar assinador: %w", err)
	}

	return stdout, 0, nil
}

func waitForServer() error {

	timeout := time.After(10 * time.Second)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout aguardando inicialização do servidor")
		case <-tick.C:
			if isServerRunning(serverPort) {
				return nil
			}
		}
	}
}

func isServerRunning(port int) bool {

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:%d/health", port)

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func checkJavaInstalled() error {

	if err := exec.Command("java", "-version").Run(); err != nil {
		if runtime.GOOS == "windows" {
			return fmt.Errorf(
				"Java não encontrado no PATH.\n→ Instale o JDK 21+ e adicione o caminho 'C:\\Program Files\\Java\\jdk-XX\\bin' nas variáveis de ambiente",
			)
		}
		return fmt.Errorf(
			"Java não encontrado no PATH.\n→ Instale o JDK 21+ e configure o PATH corretamente",
		)
	}

	return nil
}

// getAssinadorJarPath resolve o caminho do JAR na seguinte ordem de prioridade:
//  1. Variável de ambiente ASSINADOR_JAR (usada em CI e em testes de contrato)
//  2. Caminho relativo convencional do projeto
func getAssinadorJarPath() (string, error) {
	// 1. Variável de ambiente — tem precedência absoluta
	if env := os.Getenv("ASSINADOR_JAR"); env != "" {
		if _, err := os.Stat(env); err != nil {
			return "", fmt.Errorf(
				"ASSINADOR_JAR definido como %q mas arquivo não encontrado: %w",
				env, err,
			)
		}
		return env, nil
	}

	// 2. Caminho relativo convencional
	jarPath := filepath.Join("..", "assinador", "target", "assinador.jar")

	if _, err := os.Stat(jarPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf(
				"assinador.jar não encontrado em %q.\n→ Execute: cd ../assinador && mvn package",
				jarPath,
			)
		}
		return "", fmt.Errorf("erro ao verificar assinador.jar: %w", err)
	}

	return jarPath, nil
}