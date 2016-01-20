package main

import (
	"os"
	"os/exec"
	"strings"
)

type commandAndArgs struct {
	name string
	args []string
}

func parseCommand(cmd string) *commandAndArgs {
	parts := strings.Split(cmd, " ")
	if len(parts) == 1 {
		return &commandAndArgs{cmd, []string{}}
	}
	return &commandAndArgs{parts[0], parts[1:]}
}

func runCommand(c *commandAndArgs) error {
	cmd := exec.Command(c.name, c.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
