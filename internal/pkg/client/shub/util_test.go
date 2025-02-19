// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2018-2020, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package shub

import (
	"fmt"
	"os"
	"testing"

	"github.com/apptainer/apptainer/internal/pkg/test"
	useragent "github.com/apptainer/apptainer/pkg/util/user-agent"
)

var validShubURIs = []string{
	`shub://username/container`,
	`shub://username/container:tag`,
	`shub://username/container@00000000000000000000000000000000`,
	`shub://registry/username/container`,
	`shub://registry/with/levels/username/container`,
	`shub://registry/user-name/container-with-dash`,
	`shub://registry/username/container.with.period`,
	`shub://username/container:tag-with-dash`,
	`shub://username/container:tag_wtih_underscore`,
	`shub://username/container:tag.with.period`,
	`shub://myprivateregistry.sylabs.io/sylabs/container:latest`,
}

func TestMain(m *testing.M) {
	useragent.InitValue("apptainer", "v0.1.0-30-g67692d50f-dirty")

	os.Exit(m.Run())
}

// TestShubParser checks if the Shub ref parser is working as expected
func TestIsShubPullRef(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	invalidShubURIs := []string{
		`shub://username/`,
		`shub://username/container:`,
		`shub://username/container@`,
		`shub://username/container@0000000000000000000000000000000`,
		`shub://username/container@000000000000000000000000000000000`,
		`shub://username/container@abcdefghijklmnopqrstuvwxyz123456`,
		`shub://registry/user.name/container`,
		`shub://username.with.period/container:tag`,
		`shub://-username/container:`,
		`shub://username-/container:`,
		`shub://-registry/username/container:`,
		`shub://registry-/username/container:`,
	}

	for _, uri := range validShubURIs {
		t.Run(fmt.Sprintf("Valid URI: %v", uri),
			func(t *testing.T) {
				ok := isShubPullRef(uri)
				if !ok {
					t.Fatalf("failed to parse valid URI: %v", uri)
				}
			})
	}

	for _, uri := range invalidShubURIs {
		t.Run(fmt.Sprintf("Invalid URI: %v", uri),
			func(t *testing.T) {
				ok := isShubPullRef(uri)
				if ok {
					t.Fatalf("failed to parse valid URI: %v", uri)
				}
			})
	}
}

func TestShubParser(t *testing.T) {
	for _, uri := range validShubURIs {
		t.Run(fmt.Sprintf("Valid URI: %v", uri),
			func(t *testing.T) {
				sURI, err := ParseReference(uri)
				if err != nil {
					t.Fatalf("failed to parse valid URI: %v", uri)
				}
				fmt.Println(sURI.String())
			})
	}
}
