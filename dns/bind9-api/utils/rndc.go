package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// runRNDC runs rndc command and returns output
func runRNDC(args ...string) (string, error) {
	cmd := exec.Command("rndc", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return "", fmt.Errorf("rndc failed: %s", stderr.String())
	}
	return out.String(), nil
}

func ReloadBind9() (string, error) {
	return runRNDC("reload")
}
