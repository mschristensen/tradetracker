// Package main is the Leather Wallet application entrypoint.
package main

import (
	"context"
	"fmt"

	"tradetracker/internal"
	"tradetracker/internal/app/apps"
	"tradetracker/internal/app/cfg"
	"tradetracker/internal/pkg/log"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CLI command definitions.
var (
	logger logrus.FieldLogger = logrus.StandardLogger()

	rootCmd = &cobra.Command{
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}

	tradeCmd = &cobra.Command{
		Use:   "trade",
		Short: "Generates trade data.",
		RunE:  runCmd,
	}

	positionCmd = &cobra.Command{
		Use:   "position",
		Short: "Generates positions from trade data.",
		RunE:  runCmd,
	}

	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query for the position of an instrument at a given time.",
		RunE:  runCmd,
	}
)

func newApp(_ context.Context, cmd *cobra.Command, args []string) (apps.App, []string, error) {
	var err error
	var app apps.App
	switch cmd.Name() {
	case "trade":
		app, err = apps.NewTradeApp(
			cfg.DBFromEnv(),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "new trade app failed")
		}
		return app, args, nil
	case "position":
		app, err = apps.NewPositionApp(
			cfg.DBFromEnv(),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "new trade app failed")
		}
		return app, args, nil
	case "query":
		app, err = apps.NewQueryApp(
			cfg.DBFromEnv(),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "new query app failed")
		}
		return app, args, nil
	default:
		return nil, nil, fmt.Errorf("unknown command: %s", cmd.Name())
	}
}

func runCmd(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	if err := chainedCheck(
		ctx,
		envCheck,
	); err != nil {
		return errors.Wrap(err, "chained check failed")
	}
	app, args, err := newApp(cmd.Context(), cmd, args)
	if err != nil {
		return errors.Wrapf(err, "new %s app failed", cmd.Name())
	}
	return errors.Wrap(app.Run(ctx, args), "run app failed")
}

func envCheck(ctx context.Context) error {
	err := internal.ValidateEnv()
	if err != nil {
		return errors.Wrap(err, "validate env failed")
	}
	log.SetLogger(internal.LogLevel)
	return nil
}

func chainedCheck(ctx context.Context, checks ...func(context.Context) error) error {
	for _, check := range checks {
		err := check(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	err := internal.RegisterCommandFlags(rootCmd, []*internal.Flag{
		&internal.EnvFlag,
		&internal.LogLevelFlag,

		&internal.HealthPortFlag,
		&internal.PortFlag,

		&internal.MaxGoroutinesFlag,

		&internal.PostgresDatabaseFlag,
		&internal.PostgresHostFlag,
		&internal.PostgresPasswordFlag,
		&internal.PostgresPortFlag,
		&internal.PostgresUserFlag,

		&internal.MaxPGIdleConnFlag,
		&internal.MaxPGOpenConnFlag,

		&internal.RabbitMQHostFlag,
		&internal.RabbitMQUserFlag,
		&internal.RabbitMQPasswordFlag,
		&internal.RabbitMQPortFlag,
		&internal.RabbitMQVirtualHostFlag,
		&internal.RabbitMQAppIDFlag,
	})
	if err != nil {
		logger.Fatalln(err)
	}

	rootCmd.AddCommand(
		tradeCmd,
		positionCmd,
		queryCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(errors.Wrap(err, "execute root command failed"))
	}
}
