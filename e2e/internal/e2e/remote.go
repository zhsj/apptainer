// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2020-2021 Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package e2e

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/apptainer/apptainer/internal/pkg/buildcfg"
	"golang.org/x/sys/unix"
)

func SetupSystemRemoteFile(t *testing.T, testDir string) {
	Privileged(func(t *testing.T) {
		orig := filepath.Join(buildcfg.SOURCEDIR, "etc", "remote.yaml")
		dest := filepath.Join(buildcfg.APPTAINER_CONFDIR, "remote.yaml")
		source := filepath.Join(testDir, "remote.yaml")

		data, err := ioutil.ReadFile(orig)
		if err != nil {
			t.Fatalf("while reading %s: %s", orig, err)
		}
		if err := ioutil.WriteFile(source, data, 0o644); err != nil {
			t.Fatalf("while creating %s: %s", source, err)
		}
		if err := unix.Mount(source, dest, "", unix.MS_BIND, ""); err != nil {
			t.Fatalf("while mounting %s to %s: %s", source, dest, err)
		}
	})(t)
}
