package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	jdkMinVersion  = 21
	hubsaudeDir    = ".hubsaude"
	jdkDir         = "jdk"
	adoptiumAPIURL = "https://api.adoptium.net/v3/binary/latest/%d/ga/%s/%s/jdk/hotspot/normal/eclipse"
)

// resolveJava retorna o caminho do executável java a ser usado.
// Ordem de prioridade:
//  1. JDK já provisionado em ~/.hubsaude/jdk/ com versão ≥ 21
//  2. Java disponível no PATH com versão ≥ 21
//  3. Download automático via Adoptium API
func resolveJava() (string, error) {
	// 1. JDK provisionado localmente
	if javaPath, err := provisionedJavaPath(); err == nil {
		if ver, err := getJavaVersion(javaPath); err == nil && ver >= jdkMinVersion {
			slog.Debug("usando JDK provisionado localmente", "caminho", javaPath, "versão", ver)
			return javaPath, nil
		}
	}

	// 2. Java no PATH
	if javaPath, err := exec.LookPath("java"); err == nil {
		if ver, err := getJavaVersion(javaPath); err == nil && ver >= jdkMinVersion {
			slog.Debug("usando Java do PATH", "caminho", javaPath, "versão", ver)
			return javaPath, nil
		} else if err == nil {
			slog.Warn("Java encontrado no PATH mas versão incompatível",
				"versão", ver, "mínima", jdkMinVersion)
		}
	}

	// 3. Provisionar automaticamente
	slog.Info("Java 21+ não encontrado — iniciando provisionamento automático")
	return downloadAndProvisionJDK()
}

// provisionedJavaPath retorna o caminho do executável java provisionado
// em ~/.hubsaude/jdk/, ou erro se não existir.
func provisionedJavaPath() (string, error) {
	base, err := jdkBasePath()
	if err != nil {
		return "", err
	}

	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	// O JDK extraído cria um subdiretório com o nome da distribuição.
	// Buscamos o executável em qualquer subdiretório de base/jdk/.
	entries, err := os.ReadDir(base)
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		candidate := filepath.Join(base, e.Name(), "bin", javaExe)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("JDK provisionado não encontrado em %s", base)
}

// jdkBasePath retorna o diretório base de provisionamento: ~/.hubsaude/jdk/
func jdkBasePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("não foi possível determinar o diretório home: %w", err)
	}
	return filepath.Join(home, hubsaudeDir, jdkDir), nil
}

// getJavaVersion executa `java -version` e extrai o major version.
// A saída de -version vai para stderr no formato:
//
//	openjdk version "21.0.3" 2024-04-16
func getJavaVersion(javaPath string) (int, error) {
	out, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("erro ao executar %s -version: %w", javaPath, err)
	}

	return parseJavaVersion(string(out))
}

// parseJavaVersion extrai o major version de uma string de saída de `java -version`.
func parseJavaVersion(output string) (int, error) {
	// Procura padrão: "21", "21.0.3", "1.8.0_XXX"
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "version") {
			continue
		}

		// Extrai o conteúdo entre aspas: "21.0.3"
		start := strings.Index(line, `"`)
		end := strings.LastIndex(line, `"`)
		if start < 0 || end <= start {
			continue
		}

		versionStr := line[start+1 : end]
		parts := strings.Split(versionStr, ".")

		if len(parts) == 0 {
			continue
		}

		major, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		// Java 8 e anteriores reportam "1.8.0_XXX" — major real é parts[1]
		if major == 1 && len(parts) > 1 {
			minor, err := strconv.Atoi(parts[1])
			if err == nil {
				return minor, nil
			}
		}

		return major, nil
	}

	return 0, fmt.Errorf("não foi possível extrair versão de: %s", output)
}

// downloadAndProvisionJDK baixa o JDK 21 da Adoptium para a plataforma atual,
// extrai em ~/.hubsaude/jdk/ e retorna o caminho do executável java.
func downloadAndProvisionJDK() (string, error) {
	base, err := jdkBasePath()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório %s: %w", base, err)
	}

	goos, goarch := runtime.GOOS, runtime.GOARCH

	// Mapeamento para os identificadores da Adoptium API
	osName, err := adoptiumOS(goos)
	if err != nil {
		return "", err
	}

	archName, err := adoptiumArch(goarch)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(adoptiumAPIURL, jdkMinVersion, osName, archName)
	slog.Info("baixando JDK 21", "url", url, "destino", base)
	fmt.Printf("⬇️  Baixando JDK 21 para %s/%s (pode levar alguns minutos)...\n", osName, archName)

	// Determina extensão do arquivo baixado
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}

	tmpFile, err := os.CreateTemp("", "jdk-*"+ext)
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if err := downloadFile(url, tmpFile); err != nil {
		return "", fmt.Errorf(
			"falha ao baixar JDK 21\n→ verifique sua conexão com a internet\n→ instale manualmente: https://adoptium.net\n→ erro: %w",
			err,
		)
	}

	fmt.Println("📦 Extraindo JDK...")

	if goos == "windows" {
		if err := extractZip(tmpFile.Name(), base); err != nil {
			return "", fmt.Errorf("erro ao extrair JDK: %w", err)
		}
	} else {
		if err := extractTarGz(tmpFile.Name(), base); err != nil {
			return "", fmt.Errorf("erro ao extrair JDK: %w", err)
		}
	}

	javaPath, err := provisionedJavaPath()
	if err != nil {
		return "", fmt.Errorf("JDK extraído mas executável não encontrado: %w", err)
	}

	fmt.Printf("✅ JDK 21 provisionado em %s\n", base)
	slog.Info("JDK provisionado com sucesso", "caminho", javaPath)

	return javaPath, nil
}

// downloadFile baixa a URL para o arquivo destino.
func downloadFile(url string, dest *os.File) error {
	resp, err := http.Get(url) //nolint:gosec // URL construída internamente
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resposta HTTP %d de %s", resp.StatusCode, url)
	}

	_, err = io.Copy(dest, resp.Body)
	return err
}

// extractTarGz extrai um arquivo .tar.gz para o diretório destino.
func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name) //nolint:gosec

		// Proteção contra path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			return fmt.Errorf("path traversal detectado: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil { //nolint:gosec
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}

// extractZip extrai um arquivo .zip para o diretório destino (Windows).
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(dest, f.Name) //nolint:gosec

		// Proteção contra path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			return fmt.Errorf("path traversal detectado: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}

		if _, err := io.Copy(out, rc); err != nil { //nolint:gosec
			rc.Close()
			out.Close()
			return err
		}

		rc.Close()
		out.Close()
	}

	return nil
}

// adoptiumOS converte GOOS para o identificador da Adoptium API.
func adoptiumOS(goos string) (string, error) {
	switch goos {
	case "linux":
		return "linux", nil
	case "windows":
		return "windows", nil
	case "darwin":
		return "mac", nil
	default:
		return "", fmt.Errorf(
			"plataforma '%s' não suportada para provisionamento automático\n→ instale o JDK 21 manualmente: https://adoptium.net",
			goos,
		)
	}
}

// adoptiumArch converte GOARCH para o identificador da Adoptium API.
func adoptiumArch(goarch string) (string, error) {
	switch goarch {
	case "amd64":
		return "x64", nil
	case "arm64":
		return "aarch64", nil
	case "386":
		return "x86-32", nil
	default:
		return "", fmt.Errorf(
			"arquitetura '%s' não suportada para provisionamento automático\n→ instale o JDK 21 manualmente: https://adoptium.net",
			goarch,
		)
	}
}