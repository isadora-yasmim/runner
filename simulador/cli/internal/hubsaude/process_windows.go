//go:build windows

package hubsaude

import (
	"os"
)

// processAlive no Windows: os.FindProcess sempre retorna sucesso, então
// tentamos abrir o processo via Kill com sinal nulo não é possível.
// Usamos a abordagem de checar se Signal(nil)/Release falha. Como o Windows
// não suporta sinais, inferimos vivacidade tentando um Kill "soft" não é
// seguro; em vez disso, consideramos vivo se o handle pôde ser obtido e a
// porta registrada responde (a readiness real é checada à parte).
//
// Para manter o comportamento previsível e portável, retornamos true quando
// o processo pôde ser localizado; a confirmação de prontidão fica a cargo de
// isReady (actuator/health ou TCP).
func processAlive(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}

// killProcess no Windows encerra o processo diretamente (não há SIGTERM).
func killProcess(proc *os.Process) error {
	return proc.Kill()
}
