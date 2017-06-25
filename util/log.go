package util

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type Log interface {
	Debug(format string, a ...interface{})
	Output(format string, a ...interface{})
	Message(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
}

const (
	timeFormat = "2006/01/02 15:04:05"
)

type ColoredLog struct {
	l     sync.Mutex
	debug bool
}

func (c *ColoredLog) Debug(format string, a ...interface{}) {
	if c.debug {
		now := time.Now()
		c.l.Lock()
		defer c.l.Unlock()
		color.White(c.appendTime(now, format), a...)
	}
}

func (c *ColoredLog) Output(format string, a ...interface{}) {
	c.l.Lock()
	defer c.l.Unlock()
	color.White(format, a...)
}

func (ColoredLog) appendTime(stamp time.Time, str string) string {
	return fmt.Sprintf("%s [reviewbot] %s", stamp.Format(timeFormat), str)
}

func (c *ColoredLog) Message(format string, a ...interface{}) {
	now := time.Now()
	c.l.Lock()
	defer c.l.Unlock()
	color.White(c.appendTime(now, format), a...)
}

func (c *ColoredLog) Warning(format string, a ...interface{}) {
	now := time.Now()
	c.l.Lock()
	defer c.l.Unlock()
	color.Yellow(c.appendTime(now, format), a...)
}

func (c *ColoredLog) Error(format string, a ...interface{}) {
	now := time.Now()
	c.l.Lock()
	defer c.l.Unlock()
	color.Red(c.appendTime(now, format), a...)
}

func (c *ColoredLog) Info(format string, a ...interface{}) {
	now := time.Now()
	c.l.Lock()
	defer c.l.Unlock()
	color.Cyan(c.appendTime(now, format), a...)
}

type LogrusLog struct {
	log *logrus.Logger
}

func (l *LogrusLog) Debug(format string, a ...interface{}) {
	l.log.Debugf(format, a...)
}

func (l *LogrusLog) Output(format string, a ...interface{}) {
	l.log.Printf(format, a...)
}

func (l *LogrusLog) Message(format string, a ...interface{}) {
	l.log.Printf(format, a...)
}

func (l *LogrusLog) Info(format string, a ...interface{}) {
	l.log.Infof(format, a...)
}

func (l *LogrusLog) Warning(format string, a ...interface{}) {
	l.log.Warningf(format, a...)
}

func (l *LogrusLog) Error(format string, a ...interface{}) {
	l.log.Errorf(format, a...)
}

var Logger Log

func init() {
	logger := "logrus"
	debug := false
	if os.Getenv("BOT_LOGGER") != "" {
		logger = os.Getenv("BOT_LOGGER")
	}
	if len(os.Getenv("BOT_DEBUG")) > 0 {
		debug = true
	}
	switch logger {
	case "logrus":
		l := &LogrusLog{log: logrus.New()}
		Logger = l
		l.log.Out = os.Stdout

		if debug {
			l.log.Level = logrus.DebugLevel
		} else {
			l.log.Level = logrus.InfoLevel
		}
	case "color":
		c := &ColoredLog{}
		if debug {
			c.debug = true
		}
		Logger = c
	}
}
