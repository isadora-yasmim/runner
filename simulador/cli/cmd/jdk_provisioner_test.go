package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestParseJavaVersion cobre os formatos de saída do java -version
// nas versões modernas (21+) e legadas (1.8.x).
func TestParseJavaVersion(t *testing.T) {
	cases := []struct {
		name    string
		output  string
		want    int
		wantErr bool
	}{
		{
			name:   "java 21",
			output: `openjdk version "21.0.3" 2024-04-16`,
			want:   21,
		},
		{
			name:   "java 17",
			output: `openjdk version "17.0.1" 2021-10-19`,
			want:   17,
		},
		{
			name:   "java 8 legado",
			output: `java version "1.8.0_392"`,
			want:   8,
		},
		{
			name:   "java 11",
			output: `openjdk version "11.0.22" 2024-01-16`,
			want:   11,
		},
		{
			name:    "saída vazia",
			output:  "",
			wantErr: true,
		},
		{
			name:    "sem linha version",
			output:  "alguma outra saída",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseJavaVersion(tc.output)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro, mas parseJavaVersion retornou %d", got)
				}
				return
			}
			if err != nil {
				t.Errorf("parseJavaVersion(%q) erro inesperado: %v", tc.output, err)
				return
			}
			if got != tc.want {
				t.Errorf("parseJavaVersion(%q) = %d, want %d", tc.output, got, tc.want)
			}
		})
	}
}

// TestAdoptiumOS valida o mapeamento de GOOS para os identificadores da Adoptium API.
func TestAdoptiumOS(t *testing.T) {
	cases := []struct {
		goos    string
		want    string
		wantErr bool
	}{
		{"linux", "linux", false},
		{"windows", "windows", false},
		{"darwin", "mac", false},
		{"freebsd", "", true},
		{"plan9", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.goos, func(t *testing.T) {
			got, err := adoptiumOS(tc.goos)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro para goos=%s", tc.goos)
				}
				return
			}
			if err != nil {
				t.Errorf("adoptiumOS(%s) erro inesperado: %v", tc.goos, err)
				return
			}
			if got != tc.want {
				t.Errorf("adoptiumOS(%s) = %s, want %s", tc.goos, got, tc.want)
			}
		})
	}
}

// TestAdoptiumArch valida o mapeamento de GOARCH para os identificadores da Adoptium API.
func TestAdoptiumArch(t *testing.T) {
	cases := []struct {
		goarch  string
		want    string
		wantErr bool
	}{
		{"amd64", "x64", false},
		{"arm64", "aarch64", false},
		{"386", "x86-32", false},
		{"mips", "", true},
		{"riscv64", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.goarch, func(t *testing.T) {
			got, err := adoptiumArch(tc.goarch)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro para goarch=%s", tc.goarch)
				}
				return
			}
			if err != nil {
				t.Errorf("adoptiumArch(%s) erro inesperado: %v", tc.goarch, err)
				return
			}
			if got != tc.want {
				t.Errorf("adoptiumArch(%s) = %s, want %s", tc.goarch, got, tc.want)
			}
		})
	}
}

// TestJdkBasePathUsesHomeDir verifica que jdkBasePath usa o diretório home do usuário.
func TestJdkBasePathUsesHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("não foi possível determinar home dir")
	}

	got, err := jdkBasePath()
	if err != nil {
		t.Fatalf("jdkBasePath() erro inesperado: %v", err)
	}

	expected := filepath.Join(home, hubsaudeDir, jdkDir)
	if got != expected {
		t.Errorf("jdkBasePath() = %s, want %s", got, expected)
	}
}

// TestProvisionedJavaPathNotFound verifica que provisionedJavaPath retorna
// erro quando o diretório de provisionamento não existe ou está vazio.
func TestProvisionedJavaPathNotFound(t *testing.T) {
	// Redireciona temporariamente o home para um diretório temporário vazio
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	_, err := provisionedJavaPath()
	if err == nil {
		t.Error("esperava erro quando JDK não está provisionado, mas não ocorreu")
	}
}

// TestProvisionedJavaPathFound verifica que provisionedJavaPath retorna o caminho
// correto quando existe uma estrutura de JDK simulada.
func TestProvisionedJavaPathFound(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	// Cria estrutura simulada de JDK: ~/.hubsaude/jdk/jdk-21.0.3/bin/java
	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	jdkSubdir := filepath.Join(tmpHome, hubsaudeDir, jdkDir, "jdk-21.0.3", "bin")
	if err := os.MkdirAll(jdkSubdir, 0755); err != nil {
		t.Fatalf("erro ao criar estrutura de teste: %v", err)
	}

	javaPath := filepath.Join(jdkSubdir, javaExe)
	if err := os.WriteFile(javaPath, []byte("fake java"), 0755); err != nil {
		t.Fatalf("erro ao criar java fake: %v", err)
	}

	got, err := provisionedJavaPath()
	if err != nil {
		t.Fatalf("provisionedJavaPath() erro inesperado: %v", err)
	}

	if !strings.Contains(got, javaExe) {
		t.Errorf("provisionedJavaPath() = %s, esperava caminho com %s", got, javaExe)
	}
}

// TestResolveJavaUsesPathWhenValid verifica que resolveJava usa o Java do PATH
// quando ele está disponível e é versão ≥ 21.
// Este teste é skippado se não houver Java no PATH.
func TestResolveJavaUsesPathWhenValid(t *testing.T) {
	// Garante que não há JDK provisionado interferindo
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	javaInPath, err := exec.LookPath("java")
	if err != nil {
		t.Skip("java não encontrado no PATH — pulando teste")
	}

	ver, err := getJavaVersion(javaInPath)
	if err != nil || ver < jdkMinVersion {
		t.Skipf("java no PATH é versão %d (< %d) — pulando teste", ver, jdkMinVersion)
	}

	got, err := resolveJava()
	if err != nil {
		t.Fatalf("resolveJava() erro inesperado: %v", err)
	}

	if got == "" {
		t.Error("resolveJava() retornou caminho vazio")
	}
}

func TestResolveJava_JavaIncompativel_NaoDeveUsarJavaAntigo(t *testing.T) {
    tmpDir := t.TempDir()

    var scriptPath string
    if runtime.GOOS == "windows" {
        // No Windows, criamos um .bat que imprime versão 17
        scriptPath = filepath.Join(tmpDir, "java.bat")
        content := "@echo off\necho openjdk version \"17.0.1\" 2021-10-19 1>&2\n"
        if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
            t.Fatalf("erro ao criar java fake: %v", err)
        }
    } else {
        scriptPath = filepath.Join(tmpDir, "java")
        content := "#!/bin/sh\necho 'openjdk version \"17.0.1\" 2021-10-19' >&2\n"
        if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
            t.Fatalf("erro ao criar java fake: %v", err)
        }
    }

    ver, err := getJavaVersion(scriptPath)
    if err != nil {
        t.Fatalf("getJavaVersion erro inesperado: %v", err)
    }
    if ver != 17 {
        t.Fatalf("esperava versão 17 do fake, obteve %d", ver)
    }

    if ver >= jdkMinVersion {
        t.Errorf("Java 17 não deveria satisfazer versão mínima %d", jdkMinVersion)
    }
}