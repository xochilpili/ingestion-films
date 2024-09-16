package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/xochilpili/ingestion-films/internal/version"
)

/*
   Function: Constructor method for Application Logger
   @params:
   @returns: zerolog.Logger instance
*/

func New() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Str("service", "ingestion-films").Str("version", version.VERSION).Logger()
	zerolog.DefaultContextLogger = &logger
	return &logger
}
