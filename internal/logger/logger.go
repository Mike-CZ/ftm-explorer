package logger

import (
	"ftm-explorer/internal/config"
	"io"
	"strings"

	"github.com/op/go-logging"
)

// AppLogger defines extended logger with generic no-level logging option
type AppLogger struct {
	logging.Logger
}

// Printf implements default non-leveled output.
// We assume the information is low in importance if passed to this function, so we relay it to Debug level.
func (l *AppLogger) Printf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

// ModuleLogger provides a new instance of the Logger for a module.
func (l *AppLogger) ModuleLogger(module string) Logger {
	var sb strings.Builder
	sb.WriteString(l.Module)
	sb.WriteString(".")
	sb.WriteString(module)
	log := logging.MustGetLogger(sb.String())
	return &AppLogger{Logger: *log}
}

// New provides a new instance of the Logger.
func New(out io.Writer, cfg *config.Logger) *AppLogger {
	backend := logging.NewLogBackend(out, "", 0)

	fm := logging.MustStringFormatter(cfg.LogFormat)
	fmtBackend := logging.NewBackendFormatter(backend, fm)

	lvlBackend := logging.AddModuleLevel(fmtBackend)
	lvlBackend.SetLevel(cfg.LoggingLevel, "")

	logging.SetBackend(lvlBackend)
	l := logging.MustGetLogger("ftm-explorer")

	return &AppLogger{Logger: *l}
}
