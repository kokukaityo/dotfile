package main

import (
	engine "github.com/kokukaityo/dotfile/internal"
	"github.com/spf13/cobra"
)

func (a *application) statusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "コンフリクト状態を表示",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.Status(config, cmd.OutOrStdout())
		},
	}
}

func (a *application) deleteCategoryCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-category <name>",
		Short: "自動同期カテゴリを削除",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.DeleteCategory(config, args[0], cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func (a *application) gitignoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gitignore",
		Short: ".gitignoreを再生成",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := a.config()
			if err != nil {
				return err
			}
			return engine.GenerateGitignore(config)
		},
	}
}
