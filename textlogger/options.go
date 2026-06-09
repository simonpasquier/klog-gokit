package textlogger

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type ConfigOption func(*Config)

type Config struct {
	verbosity    int
	output       io.Writer
	vFlagName    string
	vmodFlagName string
	fixedTime    time.Time
	header       bool
	backtraceFn  func(skip int) (string, int)
}

func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{
		verbosity:    0,
		output:       os.Stderr,
		vFlagName:    "v",
		vmodFlagName: "vmodule",
		header:       true,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Config) Verbosity() flag.Value {
	return &verbosityValue{c: c}
}

func (c *Config) VModule() flag.Value {
	return &vmoduleValue{}
}

func (c *Config) AddFlags(fs *flag.FlagSet) {
	fs.Var(c.Verbosity(), c.vFlagName, "log level verbosity")
	fs.Var(c.VModule(), c.vmodFlagName, "per-module verbosity")
}

func VerbosityFlagName(name string) ConfigOption {
	return func(c *Config) { c.vFlagName = name }
}

func VModuleFlagName(name string) ConfigOption {
	return func(c *Config) { c.vmodFlagName = name }
}

func Verbosity(level int) ConfigOption {
	return func(c *Config) { c.verbosity = level }
}

func Output(w io.Writer) ConfigOption {
	return func(c *Config) { c.output = w }
}

func FixedTime(t time.Time) ConfigOption {
	return func(c *Config) { c.fixedTime = t }
}

func WithHeader(b bool) ConfigOption {
	return func(c *Config) { c.header = b }
}

func Backtrace(fn func(skip int) (string, int)) ConfigOption {
	return func(c *Config) { c.backtraceFn = fn }
}

type verbosityValue struct {
	c *Config
}

func (v *verbosityValue) String() string {
	if v.c == nil {
		return "0"
	}
	return fmt.Sprintf("%d", v.c.verbosity)
}

func (v *verbosityValue) Set(s string) error {
	var level int
	if _, err := fmt.Sscanf(s, "%d", &level); err != nil {
		return err
	}
	v.c.verbosity = level
	return nil
}

func (v *verbosityValue) Get() interface{} {
	return v.c.verbosity
}

type vmoduleValue struct {
	val string
}

func (v *vmoduleValue) String() string { return v.val }
func (v *vmoduleValue) Set(s string) error {
	v.val = s
	return nil
}
func (v *vmoduleValue) Get() interface{} { return v.val }
