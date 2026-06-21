package main

import (
	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

func (a *application) setupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "データリポジトリを初期設定",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.SetupRepository(config, a.hookFS, cmd.OutOrStdout())
		},
	}
}
