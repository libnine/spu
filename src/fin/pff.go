package fin

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func Pff() {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("./scripts/sh/pff.sh")
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
