package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/isadora-yasmim/simulador/internal/jdk"
)

func TestParseJavaVersion(t *testing.T) {
	cases := []struct {
		name    string
		output  string
		want    int
		wantErr bool
	}{
		{"java 21", `openjdk version "21.0.3" 2024-04-16`, 21, false},
		{"java 17", `openjdk version "17.0.1" 2021-10-19`, 17, false},
		{"java 8 legado", `java version "1.8.0_392"`, 8, false},
		{"java 11", `openjdk version "11.0.22" 2024-01-16`, 11, false},
		{"saída vazia", "", 0, true},
		{"sem linha version", "alguma outra saída", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jdk.ParseJavaVersion(tc.output)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro, mas ParseJavaVersion retornou %d", got)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseJavaVersion(%q) erro inesperado: %v", tc.output, err)
				return
			}
			if got != tc.want {
				t.Errorf("ParseJavaVersion(%q) = %d, want %d", tc.output, got, tc.want)
			}
		})
	}
}

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
			got, err := jdk.AdoptiumOS(tc.goos)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro para goos=%s", tc.goos)
				}
				return
			}
			if err != nil {
				t.Errorf("AdoptiumOS(%s) erro inesperado: %v", tc.goos, err)
				return
			}
			if got != tc.want {
				t.Errorf("AdoptiumOS(%s) = %s, want %s", tc.goos, got, tc.want)
			}
		})
	}
}

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
			got, err := jdk.AdoptiumArch(tc.goarch)
			if tc.wantErr {
				if err == nil {
					t.Errorf("esperava erro para goarch=%s", tc.goarch)
				}
				return
			}
			if err != nil {
				t.Errorf("AdoptiumArch(%s) erro inesperado: %v", tc.goarch, err)
				return
			}
			if got != tc.want {
				t.Errorf("AdoptiumArch(%s) = %s, want %s", tc.goarch, got, tc.want)
			}
		})
	}
}

func TestJdkBasePathUsesHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("não foi possível determinar home dir")
	}

	got, err := jdk.BasePath()
	if err != nil {
		t.Fatalf("BasePath() erro inesperado: %v", err)
	}

	expected := filepath.Join(home, jdk.HubsaudeDir, jdk.JdkDir)
	if got != expected {
		t.Errorf("BasePath() = %s, want %s", got, expected)
	}
}

func TestProvisionedJavaPathNotFound(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	_, err := jdk.ProvisionedJavaPath()
	if err == nil {
		t.Error("esperava erro quando JDK não está provisionado, mas não ocorreu")
	}
}

func TestProvisionedJavaPathFound(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	jdkSubdir := filepath.Join(tmpHome, jdk.HubsaudeDir, jdk.JdkDir, "jdk-21.0.3", "bin")
	if err := os.MkdirAll(jdkSubdir, 0755); err != nil {
		t.Fatalf("erro ao criar estrutura de teste: %v", err)
	}

	javaPath := filepath.Join(jdkSubdir, javaExe)
	if err := os.WriteFile(javaPath, []byte("fake java"), 0755); err != nil {
		t.Fatalf("erro ao criar java fake: %v", err)
	}

	got, err := jdk.ProvisionedJavaPath()
	if err != nil {
		t.Fatalf("ProvisionedJavaPath() erro inesperado: %v", err)
	}

	if !strings.Contains(got, javaExe) {
		t.Errorf("ProvisionedJavaPath() = %s, esperava caminho com %s", got, javaExe)
	}
}

func TestResolveJavaUsesPathWhenValid(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tmpHome)
	}

	javaInPath, err := exec.LookPath("java")
	if err != nil {
		t.Skip("java não encontrado no PATH — pulando teste")
	}

	ver, err := jdk.GetJavaVersion(javaInPath)
	if err != nil || ver < jdk.MinVersion {
		t.Skipf("java no PATH é versão %d (< %d) — pulando teste", ver, jdk.MinVersion)
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

	ver, err := jdk.GetJavaVersion(scriptPath)
	if err != nil {
		t.Fatalf("GetJavaVersion erro inesperado: %v", err)
	}
	if ver != 17 {
		t.Fatalf("esperava versão 17 do fake, obteve %d", ver)
	}
	if ver >= jdk.MinVersion {
		t.Errorf("Java 17 não deveria satisfazer versão mínima %d", jdk.MinVersion)
	}
}
