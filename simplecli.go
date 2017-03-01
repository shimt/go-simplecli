// MIT License
//
// Copyright (c) 2016 Shinichi MOTOKI
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package simplecli is simple CLI framework.
package simplecli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type application struct {
	Name string

	ProgramName string
	Arguments   []string
	Directory   string

	OS string
}

type CLI struct {
	Application application

	CommandLine *pflag.FlagSet

	ConfigFile string
	Config     *viper.Viper

	Log *logrus.Logger

	DebugMode   bool
	VerboseMode bool
}

// Initialize - initalize CLI struct.
func (c *CLI) Initialize() {
	// DebugMode & VerboseMode

	c.DebugMode = false
	c.VerboseMode = false

	// logrus

	c.initLogrus()

	// Aplication.Name

	name := filepath.Base(os.Args[0])
	ext := filepath.Ext(name)
	directory := filepath.Dir(os.Args[0])

	if ext != "" {
		name = name[:len(name)-len(ext)]
	}

	c.Application.Name = name

	// Application.ProgramName & Aplication.Arguments

	c.Application.ProgramName = os.Args[0]
	c.Application.Arguments = os.Args[1:]
	c.Application.Directory = directory

	// Application.OS
	c.Application.OS = runtime.GOOS

	// pflag & viper
	c.initPFlag()
	c.initViper()

	return
}

func (c *CLI) initLogrus() (err error) {
	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{}
	log.Out = os.Stderr
	log.Level = logrus.WarnLevel

	c.Log = log

	return
}

func (c *CLI) initViper() (err error) {
	config := viper.New()

	config.SetConfigName("." + c.Application.Name)
	config.SetEnvPrefix(c.Application.Name)

	// config path

	config.AddConfigPath(".")

	for _, name := range []string{"HOME", "USERPROFILE"} {
		if value := os.Getenv(name); value != "" {
			config.AddConfigPath(value)
		}
	}

	c.Config = config

	return
}

func (c *CLI) initPFlag() (err error) {
	commandLine := pflag.NewFlagSet(c.Application.ProgramName, pflag.ExitOnError)

	commandLine.StringVarP(&c.ConfigFile, "config", "c", "", "config file")
	commandLine.BoolVar(&c.DebugMode, "debug", false, "debug output")
	commandLine.BoolVarP(&c.VerboseMode, "verbose", "v", false, "verbose output")

	c.CommandLine = commandLine

	return
}

// Setup - Parse command line & read configuration file.
func (c *CLI) Setup() (err error) {
	c.BindSameName("debug")
	c.BindSameName("verbose")

	c.CommandLine.Parse(c.Application.Arguments)

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

	c.setupLogus()

	return
}

func (c *CLI) setupLogus() {
	if c.VerboseMode {
		c.Log.Level = logrus.InfoLevel
	}

	if c.DebugMode {
		c.Log.Level = logrus.DebugLevel
	}
}

// Exit - Exit CLI application.
func (c *CLI) Exit(code int) {
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
func (c *CLI) BindSameName(name string) {
	c.Config.BindPFlag(name, c.CommandLine.Lookup(name))
}

// NewCLI - New CLI instance.
func NewCLI() (cli *CLI) {
	cli = &CLI{}
	cli.Initialize()

	return
}

func init() {
}
