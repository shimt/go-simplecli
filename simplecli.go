package simplecli

import (
	"fmt"
	"os"
	"path"
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

func (c *CLI) Initialize() {
	// DebugMode & VerboseMode

	c.DebugMode = false
	c.VerboseMode = false

	// logrus

	c.initLogrus()

	// Aplication.Name

	name := path.Base(os.Args[0])
	ext := path.Ext(name)

	if ext != "" {
		name = name[:len(name)-len(ext)]
	}

	c.Application.Name = name

	// Application.ProgramName & Aplication.Arguments

	c.Application.ProgramName = os.Args[0]
	c.Application.Arguments = os.Args[1:]

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

func (c *CLI) Exit(code int) {
	os.Exit(code)
}

func (c *CLI) Exit1IfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.Exit(1)
	}
}

func (c *CLI) BindSameName(name string) {
	c.Config.BindPFlag(name, c.CommandLine.Lookup(name))
}

func NewCLI() (cli *CLI) {
	cli = &CLI{}
	cli.Initialize()

	return
}

func init() {
}
