package cmd

import "github.com/spf13/cobra"

// RootCommand is the root of command.
var RootCommand = &cobra.Command{
	Use:   "appsrv [command]",
	Short: "Launch your application here.",
}
