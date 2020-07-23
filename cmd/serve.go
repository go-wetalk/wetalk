package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	serve := cobra.Command{
		Use:   "serve",
		Short: "start and serve HTTP server.",
	}
	serve.AddCommand(createServeAppCommand(), createServeAdminCommand())
	RootCommand.AddCommand(&serve)
}
