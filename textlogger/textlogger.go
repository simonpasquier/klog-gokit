package textlogger

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
)

var TimeNow = time.Now

type KlogBufferWriter interface {
	WriteKlogBuffer([]byte)
}

func NewLogger(c *Config) logr.Logger {
	return logr.New(&tlogger{config: c})
}

type tlogger struct {
	callDepth int
	prefix    string
	values    []interface{}
	config    *Config
}

func (l *tlogger) Init(info logr.RuntimeInfo) {
	l.callDepth = info.CallDepth
}

func (l *tlogger) Enabled(level int) bool {
	return level <= l.config.verbosity
}

func (l *tlogger) Info(_ int, msg string, kvList ...interface{}) {
	l.print("I", msg, kvList)
}

func (l *tlogger) Error(err error, msg string, kvList ...interface{}) {
	if err != nil {
		kvList = append([]interface{}{"err", err}, kvList...)
	}
	l.print("E", msg, kvList)
}

func (l *tlogger) print(severity string, msg string, kvList []interface{}) {
	args := make([]interface{}, 0, len(l.values)+len(kvList)+4)
	if l.prefix != "" {
		args = append(args, l.prefix+": ")
	}
	args = append(args, severity, " ", msg)
	allKV := append(append([]interface{}{}, l.values...), kvList...)
	for i := 0; i+1 < len(allKV); i += 2 {
		args = append(args, fmt.Sprintf(" %s=%q", allKV[i], fmt.Sprint(allKV[i+1])))
	}
	args = append(args, "\n")
	_, _ = fmt.Fprint(l.config.output, args...)
}

func (l *tlogger) WithCallDepth(depth int) logr.LogSink {
	clone := *l
	clone.callDepth += depth
	return &clone
}

func (l *tlogger) WithName(name string) logr.LogSink {
	clone := *l
	if clone.prefix != "" {
		clone.prefix += "." + name
	} else {
		clone.prefix = name
	}
	return &clone
}

func (l *tlogger) WithValues(kvList ...interface{}) logr.LogSink {
	clone := *l
	clone.values = append(append([]interface{}{}, l.values...), kvList...)
	return &clone
}

func (l *tlogger) WriteKlogBuffer(data []byte) {
	_, _ = l.config.output.Write(data)
}

var _ logr.LogSink = &tlogger{}
var _ logr.CallDepthLogSink = &tlogger{}
var _ KlogBufferWriter = &tlogger{}
