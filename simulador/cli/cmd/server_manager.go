package cmd

import (
	"errors"
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

func waitForServer() error {

	timeout := time.After(10 * time.Second)

	// time.NewTicker + defer Stop() evita vazamento de goroutine (SA1015).
	// time.Tick() nunca libera o ticker quando a função retorna antes do tick.
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {

		case <-timeout:
			return errors.New("timeout aguardando inicialização do servidor")

		case <-ticker.C:
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

	url := fmt.Sprintf(
		"http://localhost:%d/health",
		port,
	)

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func checkJavaInstalled() error {

	command := exec.Command("java", "-version")

	err := command.Run()
	if err != nil {

		// ST1005: strings de erro não devem começar com letra maiúscula.
		if runtime.GOOS == "windows" {
			return errors.New(
				"java não encontrado no PATH.\n→ instale o JDK 21+ e adicione 'C:\\Program Files\\Java\\jdk-XX\\bin' nas variáveis de ambiente",
			)
		}

		return errors.New(
			"java não encontrado no PATH.\n→ instale o JDK 21+ e configure o PATH corretamente",
		)
	}

	return nil
}

func getAssinadorJarPath() (string, error) {

	jarPath := filepath.Join("..", "assinador", "target", "assinador-1.0-SNAPSHOT.jar")

	if _, err := os.Stat(jarPath); err != nil {
		if os.IsNotExist(err) {
			return "", errors.New(
				"assinador.jar não encontrado.\n→ execute: cd ..\\assinador && mvn clean package",
			)
		}

		return "", fmt.Errorf("erro ao verificar assinador.jar: %w", err)
	}

	return jarPath, nil
}