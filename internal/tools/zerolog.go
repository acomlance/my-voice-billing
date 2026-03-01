package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	zerologlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func ParseLevel(s string) zerolog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func InitColorfulLogger() {
	os.Setenv("TERM", "xterm-256color")

	var consoleWriter = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatLevel: func(i interface{}) string {
			if ll, ok := i.(string); ok {
				switch ll {
				case "trace":
					return "\x1b[37mTRC\x1b[0m"
				case "debug":
					return "\x1b[36mDBG\x1b[0m"
				case "info":
					return "\x1b[32mINF\x1b[0m"
				case "warn":
					return "\x1b[33mWRN\x1b[0m"
				case "error":
					return "\x1b[31mERR\x1b[0m"
				case "fatal":
					return "\x1b[35mFTL\x1b[0m"
				case "panic":
					return "\x1b[35mPNC\x1b[0m"
				default:
					return ll
				}
			}
			return "???"
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("\x1b[1m%s\x1b[0m", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("\x1b[36m%s\x1b[0m=", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("\x1b[1m%s\x1b[0m", i)
		},
	}

	var logger = zerolog.New(consoleWriter).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerologlog.Logger = logger
}

func InitFromConfig(level string, dir string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	lvl := ParseLevel(level)
	if dir == "" {
		InitColorfulLogger()
		zerolog.SetGlobalLevel(lvl)
		return
	}
	var l Logger
	l.InitLogger(dir, lvl)
	zerolog.DefaultContextLogger = &l.Logger
	zerolog.SetGlobalLevel(lvl)
	zerologlog.Logger = l.Logger
}

type Logger struct {
	zerolog.Logger
}

func (l *Logger) InitLogger(logPath string, level zerolog.Level) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var logFile string
	if stat, err := os.Stat(logPath); err == nil && stat.IsDir() {
		logFile = filepath.Join(logPath, "app.log")
	} else {
		logFile = logPath
	}
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		zerologlog.Panic().Err(err).Str("dir", dir).Msg("Failed to create log directory")
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		zerologlog.Panic().Err(err).Str("file", logFile).Msg("Failed to open log file")
	}
	out := io.MultiWriter(file, os.Stdout)
	l.Logger = zerolog.New(out).Level(level).With().Timestamp().Logger()
}
