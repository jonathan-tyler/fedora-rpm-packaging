package platform

import "io"

type Command struct {
	Name   string
	Args   []string
	Dir    string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer
}
