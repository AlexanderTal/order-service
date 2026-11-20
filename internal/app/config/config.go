package config

import (
	"io"
	"os"
	"time"

	"github.com/AlexanderTal/order-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

type LoadArgs struct {
	Output          io.Writer `json:"-"`
	EnableSimpleLog bool
}

func createLogger(args LoadArgs) zerolog.Logger {
	out := args.Output
	if out == nil {
		out = os.Stdout
	}

	var logger zerolog.Logger
	if args.EnableSimpleLog {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).
			With().Timestamp().Logger()
	} else {
		logger = zerolog.New(out).With().Timestamp().Logger()
	}

	return logger
}

func Load(args LoadArgs) {
	logger := createLogger(args)
	log.Logger = logger.Level(zerolog.DebugLevel)

	log.Debug().Msg("Logger initialized with Debug level")

	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal().Err(err).Msg("failed to process env config")
	}

	level, err := zerolog.ParseLevel(Root.Monitor.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
		log.Warn().Str("log_level", Root.Monitor.LogLevel).Msg("invalid log level, defaulting to info")
	}

	log.Logger = logger.Level(level)
	log.Info().Str("log_level", level.String()).Msg("Logger re-initialized with config level")
}
