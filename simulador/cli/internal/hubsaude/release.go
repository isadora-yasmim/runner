package hubsaude

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Constantes da release do Simulador HubSaúde publicada pelo professor.
//
// A tag segue o esquema hubsaude-validador-ui-vX.Y.Z e o asset principal é
// hubsaude-validador-ui-X.Y.Z-exec.jar (Spring Boot fat JAR). Um arquivo
// checksums.txt (formato sha256sum) acompanha cada release.
//
// Referência (commit fixo, não main, para preservar rastreabilidade):
// https://github.com/kyriosdata/runner/releases/tag/hubsaude-validador-ui-v0.1.11
const (
	githubOwner = "kyriosdata"
	githubRepo  = "runner"

	// tagPrefix identifica as releases do simulador entre outras do mesmo repo.
	tagPrefix = "hubsaude-validador-ui-"

	// assetJarSuffix casa o JAR executável independentemente da versão.
	assetJarSuffix = "-exec.jar"

	// checksumsAsset é o nome fixo do arquivo de checksums na release.
	checksumsAsset = "checksums.txt"

	// httpTimeout limita chamadas à API do GitHub (não o download do JAR).
	apiTimeout = 15 * time.Second
)

// ghAsset representa um asset de release retornado pela API do GitHub.
type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// ghRelease representa a resposta da API de releases do GitHub.
type ghRelease struct {
	TagName string    `json:"tag_name"`
	Name    string    `json:"name"`
	Assets  []ghAsset `json:"assets"`
}

// Release descreve a versão do simulador resolvida remotamente, já com os
// metadados necessários para download e verificação de integridade.
type Release struct {
	Tag         string // ex.: hubsaude-validador-ui-v0.1.11
	Version     string // ex.: 0.1.11
	JarName     string // ex.: hubsaude-validador-ui-0.1.11-exec.jar
	JarURL      string
	JarSize     int64
	ChecksumURL string
}

var versionRe = regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)

// apiBaseURL é variável para permitir testes com servidor HTTP fake.
var apiBaseURL = "https://api.github.com"

// httpClient é reutilizado; o download do JAR usa um timeout maior próprio.
func newAPIClient() *http.Client {
	return &http.Client{Timeout: apiTimeout}
}

// LatestRelease consulta a API do GitHub e resolve a release mais recente do
// simulador HubSaúde. Filtra pelo prefixo de tag para ignorar outras releases
// do mesmo repositório.
func LatestRelease() (*Release, error) {
	// /releases/latest pode apontar para uma release de outro componente do
	// mesmo repo; por isso listamos e filtramos pelo prefixo de tag.
	url := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=30", apiBaseURL, githubOwner, githubRepo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao montar requisição à API do GitHub: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := newAPIClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar releases do GitHub: %w\n→ verifique sua conexão de rede", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New(
			"limite de requisições da API do GitHub excedido (HTTP 403)\n→ defina GITHUB_TOKEN com um token pessoal para aumentar o limite",
		)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do GitHub retornou HTTP %d ao listar releases", resp.StatusCode)
	}

	var releases []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("erro ao interpretar resposta da API do GitHub: %w", err)
	}

	for _, r := range releases {
		if !strings.HasPrefix(r.TagName, tagPrefix) {
			continue
		}
		return buildRelease(r)
	}

	return nil, fmt.Errorf("nenhuma release com prefixo %q encontrada em %s/%s", tagPrefix, githubOwner, githubRepo)
}

// buildRelease converte a resposta da API no nosso tipo Release, localizando
// o asset do JAR e o de checksums.
func buildRelease(r ghRelease) (*Release, error) {
	rel := &Release{Tag: r.TagName}

	if m := versionRe.FindStringSubmatch(r.TagName); m != nil {
		rel.Version = m[1]
	}

	for _, a := range r.Assets {
		switch {
		case strings.HasSuffix(a.Name, assetJarSuffix):
			rel.JarName = a.Name
			rel.JarURL = a.BrowserDownloadURL
			rel.JarSize = a.Size
		case a.Name == checksumsAsset:
			rel.ChecksumURL = a.BrowserDownloadURL
		}
	}

	if rel.JarURL == "" {
		return nil, fmt.Errorf("release %s não contém asset com sufixo %q", r.TagName, assetJarSuffix)
	}
	return rel, nil
}

// EnsureJar garante que o JAR da release esteja disponível em cache local,
// baixando-o e verificando o checksum quando necessário. Retorna o caminho
// absoluto do JAR pronto para uso.
//
// A comparação de versão é feita pelo nome do arquivo em cache: se o JAR da
// versão remota já existe localmente, o download é pulado (US-03.4).
func EnsureJar(rel *Release) (string, error) {
	dir, err := JarCacheDir()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(dir, rel.JarName)

	if _, err := os.Stat(dest); err == nil {
		slog.Info("simulador.jar já em cache, download dispensado", "versao", rel.Version, "caminho", dest)
		return dest, nil
	}

	slog.Info("baixando simulador.jar", "versao", rel.Version, "tamanho_bytes", rel.JarSize)
	if err := downloadFile(rel.JarURL, dest); err != nil {
		return "", err
	}

	if rel.ChecksumURL != "" {
		if err := verifyChecksum(dest, rel.JarName, rel.ChecksumURL); err != nil {
			// Remove o arquivo corrompido para não envenenar o cache.
			_ = os.Remove(dest)
			return "", err
		}
		slog.Info("checksum SHA256 verificado com sucesso", "arquivo", rel.JarName)
	} else {
		slog.Warn("release sem checksums.txt; integridade não verificada", "versao", rel.Version)
	}

	return dest, nil
}

// downloadFile baixa url para dest de forma atômica (escreve em .part e
// renomeia ao final), evitando deixar arquivos truncados no cache.
func downloadFile(url, dest string) error {
	// Download de JAR grande (~164 MB): timeout generoso.
	client := &http.Client{Timeout: 10 * time.Minute}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("erro ao montar requisição de download: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao baixar %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download retornou HTTP %d para %s", resp.StatusCode, url)
	}

	tmp := dest + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return fmt.Errorf("erro ao gravar download: %w", err)
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("erro ao finalizar arquivo: %w", err)
	}

	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("erro ao mover arquivo para o cache: %w", err)
	}
	return nil
}

// verifyChecksum baixa o arquivo de checksums e confere o SHA256 de jarPath
// contra a linha correspondente a jarName. Formato esperado (sha256sum):
//
//	<hash_hex>  <nome_do_arquivo>
func verifyChecksum(jarPath, jarName, checksumURL string) error {
	expected, err := fetchExpectedChecksum(checksumURL, jarName)
	if err != nil {
		return err
	}

	actual, err := fileSHA256(jarPath)
	if err != nil {
		return err
	}

	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf(
			"checksum SHA256 não confere para %s\n→ esperado: %s\n→ obtido:   %s",
			jarName, expected, actual,
		)
	}
	return nil
}

func fetchExpectedChecksum(checksumURL, jarName string) (string, error) {
	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Get(checksumURL)
	if err != nil {
		return "", fmt.Errorf("erro ao baixar checksums.txt: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download de checksums.txt retornou HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler checksums.txt: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		// O nome pode vir com prefixo '*' (modo binário do sha256sum).
		name := strings.TrimPrefix(fields[len(fields)-1], "*")
		if name == jarName {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("checksums.txt não contém entrada para %s", jarName)
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo para checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("erro ao calcular checksum: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
