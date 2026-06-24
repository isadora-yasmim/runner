package cmd

// Este arquivo é um wrapper fino sobre internal/jdk, expondo resolveJava()
// para uso dentro do pacote cmd (assinatura CLI).
// A lógica real vive em internal/jdk/provisioner.go, compartilhada também
// pelo internal/hubsaude.

import "github.com/isadora-yasmim/simulador/internal/jdk"

// resolveJava retorna o caminho do executável java a ser usado.
// Delega para jdk.ResolveJava() que implementa a lógica completa:
//  1. JDK provisionado em ~/.hubsaude/jdk/ com versão ≥ 21
//  2. Java disponível no PATH com versão ≥ 21
//  3. Download automático via Adoptium API
func resolveJava() (string, error) {
	return jdk.ResolveJava()
}
