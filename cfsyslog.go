package cfsyslog

import (
	"fmt"
	cflog "github.com/cloudflare/golog/logger"
	golog "github.com/codemedic/go-log"
)

func toCFLevels(level golog.Level) cflog.Level {
	switch level {
	case golog.Disabled:
		return cflog.Levels.Off
	case golog.Info:
		return cflog.Levels.Info
	case golog.Warning:
		return cflog.Levels.Warn
	case golog.Error:
		return cflog.Levels.Error
	default:
		return cflog.Levels.Debug
	}
}

type cfSyslogLogger struct {
	golog.LevelledLogger
	golog.PrintLevelledLogger
	golog.StdLogSorter
	golog.SyslogTag
	logger *cflog.Logger
}

func (c *cfSyslogLogger) PrefixLogf(level golog.Level, _ int, prefix, format string, value ...interface{}) {
	if c.logger == nil {
		return
	}

	if !level.IsEnabled(c.Level()) {
		return
	}

	c.logger.Printf(toCFLevels(level), prefix, format, value...)
}

func (c *cfSyslogLogger) Write(p []byte) (n int, err error) {
	level := c.SortStdlog(c.PrintLevel(), p)
	if level.IsEnabled(c.Level()) {
		c.Logf(level, 0, "%s", string(p))
	}

	return len(p), nil
}

func (c *cfSyslogLogger) Logf(level golog.Level, _ int, format string, value ...interface{}) {
	if c.logger == nil {
		return
	}

	if !level.IsEnabled(c.Level()) {
		return
	}

	c.logger.Printf(toCFLevels(level), "", format, value...)
}

func (c *cfSyslogLogger) Close() {
	c.SetLevel(golog.Disabled)
}

// NewCFSyslog creates a new Cloudflare syslog logger with the specified options. See [cloudflare/golog](https://pkg.go.dev/github.com/cloudflare/golog/logger) for more details.
// Note that function modifies global state within the Cloudflare golog package. As a workaround, you must initialize the logger before any goroutines that use the logger are started.
func NewCFSyslog(opts ...golog.Option) (log golog.Log, err error) {
	l := &cfSyslogLogger{}

	// apply default options first
	if err = golog.SyslogDefaultOptions.Apply(l); err != nil {
		err = fmt.Errorf("error applying default options: %w", err)
		return
	}

	// apply any specified options
	for _, opt := range opts {
		if err = opt.Apply(l); err != nil {
			err = fmt.Errorf("error applying option: %w", err)
			return
		}
	}

	level := l.Level()
	if level == golog.Disabled {
		return
	}

	tag := l.GetSyslogTag()
	if tag != "" {
		if err = cflog.SetLogName(tag); err != nil {
			err = fmt.Errorf("error setting syslog tag: %w", err)
			return
		}
	}

	cfl := cflog.New(toCFLevels(level))
	if cfl == nil {
		err = fmt.Errorf("error creating syslog logger")
		return
	}

	l.logger = cfl
	log = golog.NewWithLogger(l)
	return
}
