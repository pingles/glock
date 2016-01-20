package main

import (
	"os"
	"os/exec"
)

type commandAndArgs struct {
	name string
	args []string
}

func parseCommand(cmd []string) *commandAndArgs {
	return &commandAndArgs{cmd[0], cmd[1:]}
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
