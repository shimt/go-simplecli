package simplecli

import (
	"os"
	"reflect"
	"testing"
)

func TestApplicationName(t *testing.T) {

	cli := NewCLI()

	t.Log("os.Args[0]: ", os.Args[0])
	t.Log("cli.Application.Name: ", cli.Application.Name)

	if cli.Application.Name != "go-simplecli.test" {
		t.Error("cli.Application.Name is not 'go-simplecli.test'")
	}
}

func TestApplicationProgramName(t *testing.T) {
	cli := NewCLI()

	t.Log("os.Args[0]: ", os.Args[0])
	t.Log("cli.Application.ProgramName: ", cli.Application.ProgramName)

	if os.Args[0] != cli.Application.ProgramName {
		t.Error("cli.Application.ProgramName is not `os.Args[0]`")
	}
}

func TestApplicationArguments(t *testing.T) {
	cli := NewCLI()

	t.Log("os.Args[1:]: ", os.Args[1:])
	t.Log("cli.Application.Arguments: ", cli.Application.Arguments)

	if !reflect.DeepEqual(os.Args[1:], cli.Application.Arguments) {
		t.Error("cli.Application.ProgramName is not `os.Args[1:]`")
	}
}

func TestViperConfigPath(t *testing.T) {
	// cli := NewCLI()

	t.Log("$HOME ", os.Getenv("HOME"))
	t.Log("$USERPROFILE ", os.Getenv("USERPROFILE"))
}

func TestMain(m *testing.M) {
	//start test
	code := m.Run()

	//termination
	os.Exit(code)
}
