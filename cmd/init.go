package main

import (
	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

func (a *application) initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init [path]",
		Short: "データリポジトリを新規作成",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "~/dotfiles"
			if len(args) == 1 {
				target = args[0]
			}
			return engine.InitializeRepository(target, a.templateFS, a.engineVersion, a.hookFS, cmd.OutOrStdout())
		},
	}
}
