package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ------------------------------------------------------------------ isServerRunning

func TestIsServerRunning_ServidorAtivo_DeveRetornarTrue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	port := portFromURL(t, srv.URL)
	if !isServerRunning(port) {
		t.Errorf("esperava true para servidor ativo na porta %d", port)
	}
}

func TestIsServerRunning_ServidorInativo_DeveRetornarFalse(t *testing.T) {
	port := freePort(t)
	if isServerRunning(port) {
		t.Errorf("esperava false para porta fechada %d", port)
	}
}

func TestIsServerRunning_ServidorRetorna500_DeveRetornarFalse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	port := portFromURL(t, srv.URL)
	if isServerRunning(port) {
		t.Errorf("esperava false para servidor retornando 500")
	}
}

// ------------------------------------------------------------------ resolveJava

func TestResolveJava_JavaPresente_NaoDeveRetornarErro(t *testing.T) {
	// Em CI o Java está disponível (setup-java no workflow).
	// resolveJava substitui checkJavaInstalled: detecta versão e provisiona se necessário.
	if _, err := resolveJava(); err != nil {
		t.Skipf("Java não encontrado no ambiente de teste: %v", err)
	}
}

// ------------------------------------------------------------------ getAssinadorJarPath

func TestGetAssinadorJarPath_JarAusente_DeveRetornarErroDescritivo(t *testing.T) {
	_, err := getAssinadorJarPath()
	if err == nil {
		t.Skip("JAR encontrado no path padrão — pulando teste de ausência")
	}
	msg := err.Error()
	if len(msg) == 0 {
		t.Error("mensagem de erro não pode ser vazia")
	}
	if !containsAny(msg, "assinador", "mvn", "jar") {
		t.Errorf("mensagem de erro pouco descritiva: %q", msg)
	}
}

// ------------------------------------------------------------------ waitForServer

func TestWaitForServer_ServidorJaAtivo_DeveRetornarNil(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	original := serverPort
	serverPort = portFromURL(t, srv.URL)
	defer func() { serverPort = original }()

	if err := waitForServer(); err != nil {
		t.Errorf("esperava nil, obteve: %v", err)
	}
}

// ------------------------------------------------------------------ mock health com corpo JSON correto

func TestIsServerRunning_HealthRetornaJSON_DeveRetornarTrue(t *testing.T) {
	type healthResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(healthResp{Success: true, Message: "Assinador HTTP ativo"})
	}))
	defer srv.Close()

	port := portFromURL(t, srv.URL)
	if !isServerRunning(port) {
		t.Error("esperava true para health check JSON válido")
	}
}

// ------------------------------------------------------------------ helpers

func portFromURL(t *testing.T, rawURL string) int {
	t.Helper()
	_, portStr, err := net.SplitHostPort(rawURL[len("http://"):])
	if err != nil {
		t.Fatalf("não foi possível extrair porta de %q: %v", rawURL, err)
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return port
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("não foi possível obter porta livre: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
	}
	return false
}
