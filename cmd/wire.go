//+build wireinject

package cmd

import (
	"appsrv/pkg"
	"appsrv/pkg/auth"

	"github.com/google/wire"
	"github.com/spf13/cobra"
)

func createServeAppCommand() *cobra.Command {
	wire.Build(
		pkg.ApplicationSet,
		auth.RoleGuard,
		createServerMux,
		resolveControllerSet,
		createServerRunner,
		createApplicationRunner,
	)

	return nil
}
