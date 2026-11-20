package cmd

import (
	"strings"

	"github.com/AlexanderTal/order-service/internal/app/builder"
	"github.com/urfave/cli/v2"
)

func WebServer() *cli.Command {
	return &cli.Command{
		Name:    "web-server",
		Aliases: []string{"ws"},
		Usage:   "Start HTTP server with all routes",
		Description: strings.TrimSpace(`
Initializes all dependencies (config, DB, repositories, services, handlers)
and starts the HTTP server. Graceful shutdown on SIGINT/SIGTERM.
`),
		Action:          cmdWebServer,
		HideHelpCommand: true,
	}
}

func cmdWebServer(cCtx *cli.Context) error {
	b := builder.NewBuilder(cCtx)

	b.BuildConfig().
		BuildRepoConnPostgres().
		BuildProcHttp()

	return b.Run()
}
