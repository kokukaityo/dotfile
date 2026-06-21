package main

import (
	"fmt"
	"os"

	dotfile "github.com/kokukaityo/dotfile"
)

func main() {
	if err := execute(dotfile.TemplateFS, dotfile.Version, dotfile.HookFS); err != nil {
		fmt.Fprintf(os.Stderr, "[dotfile] Error: %v\n", err)
		os.Exit(1)
	}
}
