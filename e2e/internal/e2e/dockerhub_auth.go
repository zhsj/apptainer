// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2020, Control Command Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/apptainer/apptainer/internal/pkg/util/user"
	"github.com/apptainer/apptainer/pkg/syfs"
	"oras.land/oras-go/pkg/auth"
	"oras.land/oras-go/pkg/auth/docker"
)

const dockerHub = "docker.io"

func SetupDockerHubCredentials(t *testing.T) {
	var unprivUser, privUser *user.User

	username := os.Getenv("E2E_DOCKER_USERNAME")
	pass := os.Getenv("E2E_DOCKER_PASSWORD")

	if username == "" && pass == "" {
		t.Log("No DockerHub credentials supplied, DockerHub rate limits could be hit")
		return
	}

	unprivUser = CurrentUser(t)
	writeDockerHubCredentials(t, unprivUser.Dir, username, pass)
	Privileged(func(t *testing.T) {
		privUser = CurrentUser(t)
		writeDockerHubCredentials(t, privUser.Dir, username, pass)
	})(t)
}

func writeDockerHubCredentials(t *testing.T, dir, username, pass string) {
	configPath := filepath.Join(dir, ".apptainer", syfs.DockerConfFile)
	cli, err := docker.NewClient(configPath)
	if err != nil {
		t.Fatalf("failed to get docker auth client: %v", err)
	}
	if err := cli.LoginWithOpts(
		auth.WithLoginContext(context.Background()),
		auth.WithLoginHostname(dockerHub),
		auth.WithLoginUsername(username),
		auth.WithLoginSecret(pass),
	); err != nil {
		t.Fatalf("failed to login to dockerhub: %v", err)
	}
}
