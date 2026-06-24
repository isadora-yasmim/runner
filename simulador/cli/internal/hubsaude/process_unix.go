//go:build !windows

package hubsaude

import (
	"os"
	"syscall"
)

// processAlive verifica se o processo existe enviando o sinal 0 (não entrega
// sinal, apenas checa existência/permissão), comportamento padrão em POSIX.
func processAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// killProcess solicita término gracioso (SIGTERM) ao processo.
func killProcess(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}
