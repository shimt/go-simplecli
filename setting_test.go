// Copyright 2016 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplecli

import (
	"reflect"
	"testing"
)

func TestNewCLISetting_ConfigSearchPath(t *testing.T) {
	c := NewCLI()
	s := c.NewCLISetting()

	var tests = [][]string{
		[]string{},
		[]string{"A", "B", "C"},
	}

	for _, tt := range tests {
		c.Setup(s.ConfigSearchPath(tt...))
		if !reflect.DeepEqual(c.ConfigSearchPath, tt) {
			t.Errorf("CLI.ConfigSearchPath = %v, want %v", c.ConfigSearchPath, tt)
		}
	}
}

func TestCLISetting_ConfigFile(t *testing.T) {
	c := NewCLI()
	s := c.NewCLISetting()

	var tests = []string{
		"",
		"A",
	}

	for _, tt := range tests {
		c.Setup(s.ConfigFile(tt))
		if !reflect.DeepEqual(c.ConfigFile, tt) {
			t.Errorf("CLI.ConfigFile = %v, want %v", c.ConfigFile, tt)
		}
	}
}
