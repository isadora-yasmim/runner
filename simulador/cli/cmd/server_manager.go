package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func ensureServerRunning() error {
	if isServerRunning(serverPort) {
		slog.Debug("assinador já em execução", "porta", serverPort)
		return nil
	}

	if err := checkJavaInstalled(); err != nil {
		return err
	}

	jarPath, err := getAssinadorJarPath()
	if err != nil {
		return err
	}

	slog.Warn("assinador não encontrado, iniciando automaticamente", "porta", serverPort)

	command := exec.Command("java", "-jar", jarPath, "server", "--port", fmt.Sprintf("%d", serverPort))
	if err := command.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar assinador automaticamente: %w", err)
	}

	slog.Info("assinador iniciado automaticamente", "pid", command.Process.Pid, "porta", serverPort)

	return waitForServer()
}

func waitForServer() error {
	timeout := time.After(10 * time.Second)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf(
				"timeout aguardando servidor na porta %d\n→ Verifique se a porta está disponível e se o JAR é válido",
				serverPort,
			)
		case <-tick.C:
			if isServerRunning(serverPort) {
				slog.Debug("servidor pronto", "porta", serverPort)
				return nil
			}
		}
	}
}

func isServerRunning(port int) bool {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		slog.Debug("health check sem resposta", "porta", port, "erro", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func checkJavaInstalled() error {
	if err := exec.Command("java", "-version").Run(); err != nil {
		if runtime.GOOS == "windows" {
			return fmt.Errorf(
				"Java não encontrado no PATH\n→ Instale o JDK 21+: https://adoptium.net\n→ Adicione ao PATH: C:\\Program Files\\Java\\jdk-XX\\bin",
			)
		}
		return fmt.Errorf(
			"Java não encontrado no PATH\n→ Instale o JDK 21+: sudo apt install openjdk-21-jdk\n→ Após instalar, reinicie o terminal",
		)
	}
	return nil
}

func getAssinadorJarPath() (string, error) {
	if envPath := os.Getenv("ASSINADOR_JAR"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			slog.Debug("usando JAR de ASSINADOR_JAR", "caminho", envPath)
			return envPath, nil
		}
		return "", fmt.Errorf(
			"JAR não encontrado em ASSINADOR_JAR=%s\n→ Verifique se o arquivo existe",
			envPath,
		)
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("não foi possível determinar o diretório do executável: %w", err)
	}
	execDir := filepath.Dir(execPath)

	candidates := []string{
		filepath.Join(execDir, "assinador.jar"),
		filepath.Join(execDir, "..", "assinador", "assinador.jar"),
		filepath.Join(execDir, "..", "assinador", "target", "assinador.jar"),
	}

	for _, c := range candidates {
		clean := filepath.Clean(c)
		if _, err := os.Stat(clean); err == nil {
			slog.Debug("JAR encontrado", "caminho", clean)
			return clean, nil
		}
	}

	return "", fmt.Errorf(
		"assinador.jar não encontrado\n→ Defina: export ASSINADOR_JAR=/caminho/para/assinador.jar\n→ Ou compile: cd assinador && mvn clean package -DskipTests",
	)
}