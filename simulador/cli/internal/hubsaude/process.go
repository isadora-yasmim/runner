package hubsaude

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/isadora-yasmim/simulador/internal/jdk"
)

// DefaultPort é a porta padrão do Simulador HubSaúde (Spring Boot).
// Registrada como decisão em documentacao/adr/ (porta padrão).
const DefaultPort = 9595

// procInfo é persistido em ~/.hubsaude/simulador.pid para permitir que
// comandos subsequentes (stop/status) encontrem o processo iniciado antes.
type procInfo struct {
	PID  int `json:"pid"`
	Port int `json:"port"`
}

// PortAvailable verifica se a porta TCP está livre para bind local.
// Usado antes de iniciar o simulador (US-03.1: verificar portas).
func PortAvailable(port int) bool {
	addr := net.JoinHostPort("localhost", strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// Start inicia o simulador.jar como processo em background na porta indicada,
// registrando PID e porta em ~/.hubsaude/simulador.pid.
//
// Usa jdk.ResolveJava() para garantir Java 21+ disponível, provisionando
// automaticamente se necessário (US-04.1).
//
// Retorna erro de usuário (exit 1) para porta ocupada e erro de sistema
// (exit 2) para JVM ausente ou falha ao iniciar — a separação é feita pelo
// chamador via os código de saída.
func Start(jarPath string, port int) (int, error) {
	if running, info := isRegisteredProcessAlive(); running {
		return info.PID, fmt.Errorf("simulador já está em execução (PID %d, porta %d)", info.PID, info.Port)
	}

	if !PortAvailable(port) {
		return 0, fmt.Errorf(
			"a porta %d já está em uso por outro processo\n→ use --port <outra> para escolher uma porta livre",
			port,
		)
	}

	// Garante Java 21+ disponível, provisionando automaticamente se necessário.
	javaPath, err := jdk.ResolveJava()
	if err != nil {
		return 0, err
	}

	// O JAR é Spring Boot: a porta é configurada via --server.port.
	cmd := exec.Command(javaPath, "-jar", jarPath, "--server.port="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("erro ao iniciar simulador.jar: %w", err)
	}

	info := procInfo{PID: cmd.Process.Pid, Port: port}
	if err := writePidFile(info); err != nil {
		slog.Warn("simulador iniciado mas não foi possível registrar PID", "erro", err)
	}

	// Libera o processo filho para seguir vivo após o término do CLI.
	_ = cmd.Process.Release()

	return info.PID, nil
}

// Stop encerra o processo do simulador registrado em ~/.hubsaude/simulador.pid.
func Stop() (int, error) {
	info, err := readPidFile()
	if err != nil {
		return 0, errors.New("nenhum simulador registrado como em execução\n→ nada a encerrar")
	}

	proc, err := os.FindProcess(info.PID)
	if err != nil {
		_ = removePidFile()
		return 0, fmt.Errorf("processo PID %d não encontrado: %w", info.PID, err)
	}

	if err := killProcess(proc); err != nil {
		return 0, fmt.Errorf("erro ao encerrar simulador (PID %d): %w", info.PID, err)
	}

	_ = removePidFile()
	return info.PID, nil
}

// Status representa o estado observável do simulador.
type Status struct {
	Registered bool // há PID registrado
	PID        int
	Port       int
	Alive      bool // processo existe (liveness: processo subiu)
	Ready      bool // /actuator/health respondeu UP (readiness)
}

// Query inspeciona o estado do simulador combinando: existência do processo
// (liveness) e prontidão para requisições (readiness via /actuator/health,
// com fallback para porta TCP aberta). Não confunde "subiu" com "pronto".
func Query() Status {
	var st Status

	info, err := readPidFile()
	if err != nil {
		return st
	}
	st.Registered = true
	st.PID = info.PID
	st.Port = info.Port
	st.Alive = processAlive(info.PID)

	if st.Alive {
		st.Ready = isReady(info.Port)
	}
	return st
}

// isReady tenta o endpoint de readiness do Spring Boot Actuator; se o actuator
// não estiver exposto, faz fallback para "porta TCP aceitando conexões".
func isReady(port int) bool {
	if actuatorHealthUP(port) {
		return true
	}
	return tcpReachable(port)
}

// actuatorHealthUP consulta /actuator/health e considera pronto quando o JSON
// retorna {"status":"UP"} com HTTP 200.
func actuatorHealthUP(port int) bool {
	client := http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://localhost:%d/actuator/health", port)

	resp, err := client.Get(url)
	if err != nil {
		slog.Debug("actuator/health sem resposta", "porta", port, "erro", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return true
	}
	return strings.EqualFold(payload.Status, "UP")
}

func tcpReachable(port int) bool {
	addr := net.JoinHostPort("localhost", strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// isRegisteredProcessAlive informa se há um processo registrado e ainda vivo.
func isRegisteredProcessAlive() (bool, procInfo) {
	info, err := readPidFile()
	if err != nil {
		return false, procInfo{}
	}
	return processAlive(info.PID), info
}

// --- persistência do PID ---

func writePidFile(info procInfo) error {
	path, err := pidFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func readPidFile() (procInfo, error) {
	var info procInfo
	path, err := pidFilePath()
	if err != nil {
		return info, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return info, err
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return info, err
	}
	if info.PID == 0 {
		return info, errors.New("registro de PID inválido")
	}
	return info, nil
}

func removePidFile() error {
	path, err := pidFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
