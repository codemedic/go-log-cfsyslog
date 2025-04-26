package cfsyslog

import (
	golog "github.com/codemedic/go-log"
)

func ExampleNewCFSyslog() {
	l := golog.Must(NewCFSyslog(
		golog.OptionsMust(golog.Options(
			// set the log-level dynamically from the environment
			golog.WithLevelFromEnv("LOG_LEVEL", golog.Info),
			// set the syslog tag
			golog.WithSyslogTag("test-syslog"),
		))))

	defer l.Close()
}
