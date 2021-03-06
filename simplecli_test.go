// Copyright 2016 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplecli

import (
	"os"
	"reflect"
	"testing"
)

func TestApplicationName(t *testing.T) {

	cli := NewCLI()

	t.Log("os.Args[0]:", os.Args[0])
	t.Log("cli.Application.Name:", cli.Application.Name)

	if cli.Application.Name != "go-simplecli" {
		t.Error("cli.Application.Name is not 'go-simplecli'")
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

func TestAConfigSearchPath(t *testing.T) {
	cli := NewCLI()

	t.Log("cli.ConfigSearchPath:", cli.ConfigSearchPath)

	csp := cli.ConfigSearchPath
	i := 0

	if csp[i] != "." {
		t.Errorf("cli.ConfigSearchPath[%d] is not '.'", i)
	}
	i++

	if value := os.Getenv("HOME"); value != "" {
		t.Log("HOME:", value)
		if csp[i] != value {
			t.Errorf("cli.ConfigSearchPath[%d] is not '%s'", i, value)
		}
		i++
	}

	if value := os.Getenv("USERPROFILE"); value != "" {
		t.Log("USERPROFILE:", value)
		if csp[i] != value {
			t.Errorf("cli.ConfigSearchPath[%d] is not '%s'", i, value)
		}
		// i++
	}
}

func TestBetweenRune(t *testing.T) {
	var betweenRuneTests = []struct {
		char     rune
		lower    rune
		upper    rune
		expected bool
	}{
		{'0', '0', '9', true},
		{'5', '0', '9', true},
		{'9', '0', '9', true},
		{'A', '0', '9', false},
		{'A', 'A', 'Z', true},
		{'M', 'A', 'Z', true},
		{'Z', 'A', 'Z', true},
		{'0', 'A', 'Z', false},
	}

	for _, tt := range betweenRuneTests {
		actual := betweenRune(tt.char, tt.lower, tt.upper)
		if actual != tt.expected {
			t.Errorf(
				"betweenRuneTests('%s', '%s', '%s'): expected %t, actual %t",
				string(tt.char), string(tt.lower), string(tt.upper), tt.expected, actual)
		}
	}

}

func TestNormalizeEnvName(t *testing.T) {
	var normalizeEnvNameTests = []struct {
		input    string
		expected string
	}{
		{"test", "TEST"},
		{"test-test", "TEST_TEST"},
		{"0test-test", "_TEST_TEST"},
		{"!\"#$%&'()=-^\\", "_____________"},
		{"あいうえお", "_____"},
	}

	for _, tt := range normalizeEnvNameTests {
		actual := normalizeEnvName(tt.input)
		if actual != tt.expected {
			t.Errorf(
				"normalizeEnvName(\"%s\"): expected %s, actual %s",
				tt.input, tt.expected, actual)
		}
	}

}

func TestViperConfigPath(t *testing.T) {
	// cli := NewCLI()

	t.Log("HOME:", os.Getenv("HOME"))
	t.Log("USERPROFILE:", os.Getenv("USERPROFILE"))
}

func TestMain(m *testing.M) {
	//start test
	code := m.Run()

	//termination
	os.Exit(code)
}
