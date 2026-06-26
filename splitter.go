package main

import "strings"

func SplitArgs(s string) []string {
	var args []string
	var current strings.Builder
	inquotes := false

	for _, r := range s {
		switch {
		case r == '"' || r == '\'':
			inquotes = !inquotes
		case r == ' ' && !inquotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
