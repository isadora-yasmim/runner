package hubsaude

import (
	"fmt"
	"os"
	"path/filepath"
)

// baseDirName é o diretório gerenciado pelo Sistema Runner no HOME do usuário,
// conforme definido nas histórias US-03/US-04 (~/.hubsaude/).
const baseDirName = ".hubsaude"

// BaseDir retorna o caminho de ~/.hubsaude, criando-o se necessário.
func BaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao localizar diretório HOME do usuário: %w", err)
	}
	dir := filepath.Join(home, baseDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
	}
	return dir, nil
}

// JarCacheDir retorna ~/.hubsaude/simulador, onde os JARs do simulador são
// armazenados em cache entre execuções (US-03.4).
func JarCacheDir() (string, error) {
	base, err := BaseDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "simulador")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de cache %s: %w", dir, err)
	}
	return dir, nil
}

// pidFilePath retorna o caminho do arquivo que registra o processo do
// simulador em execução, permitindo stop/status entre invocações do CLI.
func pidFilePath() (string, error) {
	base, err := BaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "simulador.pid"), nil
}
