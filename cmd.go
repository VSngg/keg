// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package example provides the Bonzai command branch of the same name.
package keg

import (
	Z "github.com/rwxrob/bonzai/z"
)

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{

	Name:      `keg`,
	Summary:   `Knowledge Exchange Grid (KEG)`,
	Version:   `v0.0.0`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		//		help.Cmd, conf.Cmd, vars.Cmd,
	},

	Description: `
		The **{{.Name}}** command composes together the branches and
		commands used to search, read, create, and share knowledge on the
		free, decentralized, protocol-agnostic, world-wide, Knowledge
		Exchange Grid, a modern replacement for the very broken WorldWideWeb
		(see keg.pub for more).`,
}
