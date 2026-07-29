package main

import (
	"encoding/json"
	"os/exec"
)

type DistList struct {
	Os     string `json:"GOOS"`
	Arch   string `json:"GOARCH"`
	Cgo    bool   `json:"CgoSupported"`
	Fclass bool   `json:"FirstClass"`
}

func List() ([]DistList, error) {
	cmd := exec.Command("go", "tool", "dist", "list", "-json")
	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var list []DistList
	err = json.Unmarshal(o, &list)
	if err != nil {
		return []DistList{}, err
	}

	return list, nil
}
