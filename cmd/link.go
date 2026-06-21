package main

import (
	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

func (a *application) linkCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "link",
		Short: "symlinkを配置",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.LinkAll(config, cmd.OutOrStdout())
		},
	}
}
