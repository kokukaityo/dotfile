package main

import (
	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

func (a *application) pullCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "リモートから同期",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.Pull(config, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func (a *application) pushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "変更をcommitしてpush",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.Push(config, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}
