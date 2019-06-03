package logbook

import (
	"fmt"
	"log"
)

//Severity Urgency of the log message
type Severity int

//LogBook An operation logger
type LogBook struct {
	threshold  Severity
	warnings   bool
	logChannel chan *logBookItem
}

type logBookItem struct {
	Level   Severity
	Message string
}

const (
	//TRACE log level, verbosity = 3
	TRACE Severity = iota
	//DEBUG log level = verbosity = 2
	DEBUG Severity = iota
	//INFO log level, verbosity = 1
	INFO Severity = iota
	//WARN log level, verbosity = 0
	WARN Severity = iota
	//ERROR log level
	ERROR Severity = iota
	//FATAL log level
	FATAL Severity = iota
)

func severityString(lvl Severity) string {
	switch lvl {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		panic(lvl)
	}
}

func logLevel(verbosity int) Severity {
	switch verbosity {
	case 0:
		return WARN
	case 1:
		return INFO
	case 2:
		return DEBUG
	case 3:
		return TRACE
	default:
		log.Printf("ERROR: verbosity=%d but shall be 0..3\n", verbosity)
		return WARN
	}
}

//RunLogBook A routine to log messages. verbosity=3 equals to TRACE, verbosity=0 equals to WARN
func (logbook *LogBook) runLogBook() {
	for item := range logbook.logChannel {
		if item.Level == WARN && !logbook.warnings {
			continue
		}

		if item.Level >= logbook.threshold {
			log.Println(severityString(item.Level) + " " + item.Message)
		}
	}
}

//NewLogBook return a new logbook channel to log operations into
func NewLogBook(verbosity int, quietness int) *LogBook {
	//TODO log verbosity
	logbook := &LogBook{threshold: logLevel(verbosity), warnings: quietness == 0, logChannel: make(chan *logBookItem, 100)}

	go logbook.runLogBook()
	logbook.Log(INFO, fmt.Sprintf("Setting log level to %s", severityString(logbook.threshold)))
	return logbook
}

//Log append a new log item into the logbook
func (logbook *LogBook) Log(level Severity, message string) {
	logbook.logChannel <- &logBookItem{level, message}
}
