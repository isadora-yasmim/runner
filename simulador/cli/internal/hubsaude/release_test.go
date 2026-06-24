package hubsaude

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildRelease_ResolveJarEChecksum(t *testing.T) {
	r := ghRelease{
		TagName: "hubsaude-validador-ui-v0.1.11",
		Assets: []ghAsset{
			{Name: "checksums.txt", BrowserDownloadURL: "http://x/checksums.txt"},
			{Name: "hubsaude-validador-ui-0.1.11-exec.jar", BrowserDownloadURL: "http://x/jar", Size: 123},
		},
	}

	rel, err := buildRelease(r)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if rel.Version != "0.1.11" {
		t.Errorf("versão esperada 0.1.11, obtida %q", rel.Version)
	}
	if rel.JarName != "hubsaude-validador-ui-0.1.11-exec.jar" {
		t.Errorf("JarName incorreto: %q", rel.JarName)
	}
	if rel.JarSize != 123 {
		t.Errorf("JarSize incorreto: %d", rel.JarSize)
	}
	if rel.ChecksumURL == "" {
		t.Error("ChecksumURL não resolvido")
	}
}

func TestBuildRelease_SemJar_Falha(t *testing.T) {
	r := ghRelease{
		TagName: "hubsaude-validador-ui-v0.1.11",
		Assets:  []ghAsset{{Name: "checksums.txt"}},
	}
	if _, err := buildRelease(r); err == nil {
		t.Fatal("esperava erro quando não há asset -exec.jar")
	}
}

func TestFetchExpectedChecksum_FormatoSha256sum(t *testing.T) {
	const jarName = "hubsaude-validador-ui-0.1.11-exec.jar"
	const hash = "e0766817b68b414438a30208f5abc0000000000000000000000000000000000"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%s  %s\n", hash, jarName)
	}))
	defer srv.Close()

	got, err := fetchExpectedChecksum(srv.URL, jarName)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got != hash {
		t.Errorf("hash esperado %q, obtido %q", hash, got)
	}
}

func TestFetchExpectedChecksum_EntradaAusente(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "deadbeef  outro-arquivo.jar")
	}))
	defer srv.Close()

	if _, err := fetchExpectedChecksum(srv.URL, "ausente.jar"); err == nil {
		t.Fatal("esperava erro para entrada ausente no checksums.txt")
	}
}

func TestFileSHA256(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(p, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	// SHA256("hello")
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got, err := fileSHA256(p)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("sha256 esperado %s, obtido %s", want, got)
	}
}

func TestLatestRelease_FiltraPorPrefixoDeTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Primeira release é de outro componente; deve ser ignorada.
		fmt.Fprint(w, `[
		  {"tag_name":"assinatura-v9.9.9","assets":[]},
		  {"tag_name":"hubsaude-validador-ui-v0.1.11","assets":[
		     {"name":"hubsaude-validador-ui-0.1.11-exec.jar","browser_download_url":"http://x/jar","size":1},
		     {"name":"checksums.txt","browser_download_url":"http://x/c"}
		  ]}
		]`)
	}))
	defer srv.Close()

	old := apiBaseURL
	apiBaseURL = srv.URL
	defer func() { apiBaseURL = old }()

	rel, err := LatestRelease()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !strings.HasPrefix(rel.Tag, tagPrefix) {
		t.Errorf("tag resolvida não tem prefixo esperado: %q", rel.Tag)
	}
}

func TestPortAvailable(t *testing.T) {
	// Abre uma porta e confirma que PortAvailable a reporta como ocupada.
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	if PortAvailable(port) {
		t.Errorf("porta %d deveria estar ocupada", port)
	}
}

func TestCacheDirs_CriamSubdiretorios(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // Windows

	dir, err := JarCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("diretório de cache não criado: %v", err)
	}
}
