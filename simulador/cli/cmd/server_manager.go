package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
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

	fmt.Printf("⚠️ Assinador não encontrado na porta %d.\n", serverPort)
	fmt.Println("🚀 Iniciando automaticamente...")

	jarPath := "..\\assinador\\target\\assinador-1.0-SNAPSHOT.jar"

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

	fmt.Printf("✅ Assinador iniciado automaticamente.\n")
	fmt.Printf("PID: %d\n", command.Process.Pid)

	return waitForServer()
}

func waitForServer() error {

	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {

		case <-timeout:
			return fmt.Errorf("timeout aguardando inicialização do servidor")

		case <-tick:
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
