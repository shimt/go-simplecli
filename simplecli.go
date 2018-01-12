// Copyright 2016 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplecli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkg/profile"
	"github.com/shimt/go-logif"
	"github.com/shimt/go-logif/gologif"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type application struct {
	Name string

	ProgramName string
	Arguments   []string
	Directory   string

	OS   string
	Arch string
}

// CLI - CLI main structure
type CLI struct {
	Application application

	CommandLine *pflag.FlagSet

	ConfigSearchPath []string
	ConfigFile       string
	Config           *viper.Viper

	Log logif.Logger

	DebugMode   bool
	VerboseMode bool
	ProfileMode string

	InitializeError error

	profiler interface {
		Stop()
	}
}

// Initialize - initialize CLI struct.
func (c *CLI) Initialize() (err error) {
	arg0 := os.Args[0]

	// DebugMode & VerboseMode

	c.DebugMode = false
	c.VerboseMode = false

	// Logger

	if e := c.initLogger(); e != nil {
		err = errors.Wrap(e, "CLI.initLogrus")
		return
	}

	// Aplication.Name

	name := filepath.Base(arg0)
	ext := filepath.Ext(name)
	directory := filepath.Dir(arg0)

	if ext != "" {
		name = name[:len(name)-len(ext)]
	}

	c.Application.Name = name

	// Application.ProgramName & Application.Arguments

	c.Application.ProgramName = arg0
	c.Application.Arguments = os.Args[1:]
	c.Application.Directory = directory

	// Application.OS & Application.Arch
	c.Application.OS = runtime.GOOS
	c.Application.Arch = runtime.GOARCH

	// ConfigSearchPath
	configSearchPath := []string{"."}

	for _, name := range []string{"HOME", "USERPROFILE"} {
		if value := os.Getenv(name); value != "" {
			configSearchPath = append(configSearchPath, value)
		}
	}

	c.ConfigSearchPath = configSearchPath

	// pflag & viper

	if e := c.initPFlag(); e != nil {
		err = errors.Wrap(e, "CLI.initPFlag")
		return
	}

	if e := c.initViper(); e != nil {
		err = errors.Wrap(e, "CLI.initViper")
		return
	}

	return
}

func (c *CLI) initLogger() (err error) {
	c.Log = gologif.New(os.Stderr, "", log.LstdFlags)

	return
}

func betweenRune(c rune, l rune, u rune) bool {
	return l <= c && c <= u
}

func normalizeEnvName(name string) string {
	buf := []rune(strings.ToUpper(name))

	if len(buf) > 0 {
		for i, c := range buf {
			if !(betweenRune(c, 'A', 'Z') || betweenRune(c, '0', '9') || c == '_') {
				buf[i] = '_'
			}
		}

		if betweenRune(buf[0], '0', '9') {
			buf[0] = '_'
		}
	}

	return string(buf)
}

func (c *CLI) initViper() (err error) {
	config := viper.New()

	config.SetConfigName("." + c.Application.Name)
	config.SetEnvPrefix(normalizeEnvName(c.Application.Name))

	// config path
	for _, path := range c.ConfigSearchPath {
		config.AddConfigPath(path)
	}

	c.Config = config

	return
}

func (c *CLI) initPFlag() (err error) {
	commandLine := pflag.NewFlagSet(c.Application.ProgramName, pflag.ExitOnError)

	commandLine.StringVar(&c.ConfigFile, "config", "", "config file")
	commandLine.BoolVar(&c.DebugMode, "debug", false, "debug output")
	commandLine.BoolVar(&c.VerboseMode, "verbose", false, "verbose output")
	commandLine.StringVar(&c.ProfileMode, "profile", "", "profile mode (cpu/memory/mutex/block/trace)")

	c.CommandLine = commandLine

	return
}

var profileMap = map[string]func(*profile.Profile){
	"cpu":    profile.CPUProfile,
	"memory": profile.MemProfile,
	"mutex":  profile.MutexProfile,
	"block":  profile.BlockProfile,
	"trace":  profile.TraceProfile,
}

// Setup - Parse command line & read configuration file.
func (c *CLI) Setup(setups ...func()) (err error) {
	for _, f := range setups {
		f()
	}

	if e := c.BindSameName("debug"); e != nil {
		err = errors.Wrap(e, "CLI.BindSameName(\"debug\")")
		return
	}

	if e := c.BindSameName("verbose"); e != nil {
		err = errors.Wrap(e, "CLI.BindSameName(\"debug\")")
		return
	}

	if e := c.CommandLine.Parse(c.Application.Arguments); e != nil {
		err = errors.Wrap(e, "CLI.CommandLine.Parse")
		return
	}

	if c.ConfigFile != "" {
		c.Config.SetConfigFile(c.ConfigFile)
	}

	c.Config.AutomaticEnv()

	err = c.Config.ReadInConfig()
	configFileNotFoundError := false

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = nil
		configFileNotFoundError = true
	}

	if err != nil {
		return errors.Wrap(err, "failed to read configuration file")
	}

	if !configFileNotFoundError {
		c.ConfigFile = c.Config.ConfigFileUsed()
	} else {
		c.Log.Info("configuration file not found")
	}

	c.setupLogger()

	if err = c.checkProfileMode(); err != nil {
		return
	}

	return
}

func (c *CLI) setupLogger() {
	if m, ok := c.Log.(logif.LeveledLoggerModifier); ok {
		if c.VerboseMode {
			m.SetOutputLevel(logif.INFO)
		}

		if c.DebugMode {
			m.SetOutputLevel(logif.DEBUG)
		}
	}
}

func (c *CLI) checkProfileMode() (err error) {
	if c.ProfileMode != "" {
		if _, ok := profileMap[c.ProfileMode]; !ok {
			err = errors.Errorf("unknown profiler (%s)", c.ProfileMode)
			return
		}
	}

	return
}

// StartProfile - Start profiling
func (c *CLI) StartProfile() {
	if pm, ok := profileMap[c.ProfileMode]; ok {
		c.profiler = profile.Start(pm)
	}
}

// StopProfile - Stop profiling
func (c *CLI) StopProfile() {
	if c.profiler != nil {
		c.profiler.Stop()
		c.profiler = nil
	}
}

// Exit - Exit CLI application.
func (c *CLI) Exit(code int) {
	c.StopProfile()
	os.Exit(code)
}

// Exit1IfError - Exit CLI application if error.
func (c *CLI) Exit1IfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.Exit(1)
	}
}

// BindSameName - Bind viper & pflag parameter.
func (c *CLI) BindSameName(names ...string) (err error) {
	for _, name := range names {
		if e := c.Config.BindPFlag(name, c.CommandLine.Lookup(name)); e != nil {
			err = errors.Wrap(e, "CLI.BindSameName")
			return
		}
	}

	return
}

// NewCLISetting - New CLISetting instance.
func (c *CLI) NewCLISetting() CLISetting {
	return CLISetting{c}
}

// NewCLI - New CLI instance.
func NewCLI() (cli *CLI) {
	cli = &CLI{}
	cli.InitializeError = cli.Initialize()

	return
}

func init() {
}
