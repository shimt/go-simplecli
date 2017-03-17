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

	t.Log("$HOME ", os.Getenv("HOME"))
	t.Log("$USERPROFILE ", os.Getenv("USERPROFILE"))
}

func TestMain(m *testing.M) {
	//start test
	code := m.Run()

	//termination
	os.Exit(code)
}
