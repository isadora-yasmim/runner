package executor

import (
	"os"
	"os/exec"
)

func RunJar(jarPath string, args ...string) error {
	cmdArgs := append([]string{"-jar", jarPath}, args...)

	cmd := exec.Command("java", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}