package main

import (
	"fmt"
	"os"
	"os/exec"
)

func Clibuild(name string, cgoenabled string) {
	args := []string{"build", "-ldflags", "-s -w", "-trimpath"}

	if name != "" {
		args = append(args, "-o", name)
	}

	cmd := exec.Command("go", args...)

	if cgoenabled == "1" {
		cmd.Env = append(os.Environ(),
			"CGO_ENABLED=1")
	} else {
		cmd.Env = append(os.Environ(),
			"CGO_ENABLED=0")
	}

	o, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(o))
		panic(err)
	}
	fmt.Println(string(o))
}
