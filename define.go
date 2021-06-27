package logger

import (
	"github.com/ivpusic/grpool"
	"log"
)

type Logger struct {
	level int
	std   *log.Logger
	err   *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	p     *grpool.Pool
	depth int
}

const (
	LevelError = iota
	LevelWarning
	LevelInformational
	LevelDebug
)

var DefaultLogger = GetDefaultLogger(512, log.LstdFlags|log.Lmicroseconds, 1, 100, 1)