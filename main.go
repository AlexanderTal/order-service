package main

import (
	"fmt"
	"os"

	"github.com/AlexanderTal/order-service/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "order-service",
		Version: "1.0.0",
		Usage:   "Order management service",
		Commands: []*cli.Command{
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "Enable simple log format instead of JSON",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
