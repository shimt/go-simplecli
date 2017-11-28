// Copyright 2016 Shinichi MOTOKI. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package simplecli is simple CLI framework.
package simplecli

// CLISetting - CLI Setting Structure
type CLISetting struct {
	cli *CLI
}

// ConfigSearchPath - Set ConfigSearchPath
func (c *CLISetting) ConfigSearchPath(paths ...string) func() {
	return func() {
		c.cli.ConfigSearchPath = paths[:]
	}
}

// ConfigFile - Set ConfigFIle
func (c *CLISetting) ConfigFile(path string) func() {
	return func() {
		c.cli.ConfigFile = path
	}
}
