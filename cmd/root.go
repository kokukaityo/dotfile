package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

type application struct {
	templateFS    fs.FS
	hookFS        fs.FS
	engineVersion string
}

func execute(templateFS fs.FS, engineVersion string, hookFS fs.FS) error {
	app := &application{
		templateFS:    templateFS,
		hookFS:        hookFS,
		engineVersion: strings.TrimSpace(engineVersion),
	}
	return app.rootCommand().Execute()
}

func (a *application) rootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "dotfile",
		Short:         "dotfiles同期エンジン",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(
		a.initCommand(),
		a.setupCommand(),
		a.linkCommand(),
		a.pullCommand(),
		a.pushCommand(),
		a.deleteCategoryCommand(),
		a.gitignoreCommand(),
		a.statusCommand(),
		a.versionCommand(),
	)
	return root
}

func (a *application) config() (*engine.Config, error) {
	config, err := engine.Resolve(a.engineVersion)
	if err != nil {
		return nil, err
	}
	if config.VersionMismatch() {
		_, _ = fmt.Fprintf(os.Stderr, "[dotfile] WARNING: バージョン不整合 (engine=%s, data=%s)\n", config.EngineVersion, config.DataVersion)
	}
	return config, nil
}
