// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package exec

import (
	"testing"

	"gotest.tools/v3/icmd"
)

type Cmd struct {
	path string
	args []string
}

type Result struct {
	*icmd.Result
}

func Command(path string, args ...string) *Cmd {
	return &Cmd{path: path, args: args}
}

func (c *Cmd) Run(t *testing.T) *Result {
	t.Logf("Running cmd %q with args %q", c.path, c.args)

	result := icmd.RunCommand(c.path, c.args...)

	return &Result{result}
}
