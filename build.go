package main

import (
	"os"
	"os/exec"
)

func Compile(goos string, goarch string, cgosupport string, custargs string) (string, error) {
	args := []string{"build"}

	if custargs != "" {
		args = append(args, SplitArgs(custargs)...)
	}

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
		"CGO_ENABLED="+cgosupport,
	)

	o, err := cmd.CombinedOutput()
	if err != nil {
		return string(o), err
	}

	return string(o), nil
}
