// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2021, Sylabs Inc. All rights reserved.
// Copyright (c) 2020, Control Command Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package apptainer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/apptainer/apptainer/internal/pkg/remote"
	"github.com/apptainer/apptainer/internal/pkg/remote/endpoint"
)

func RemoteRemoveKeyserver(name, uri string) error {
	// Explicit handling of corner cases: name and uri must be valid strings
	if strings.TrimSpace(uri) == "" {
		return fmt.Errorf("invalid URI: cannot have empty URI")
	}

	// opening config file
	file, err := os.OpenFile(remote.SystemConfigPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("while opening remote config file: %s", err)
	}
	defer file.Close()

	// read file contents to config struct
	c, err := remote.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("while parsing remote config data: %s", err)
	}

	var ep *endpoint.Config

	if name == "" {
		ep, err = c.GetDefault()
	} else {
		ep, err = c.GetRemote(name)
	}

	if err != nil {
		return fmt.Errorf("no endpoint found: %s", err)
	} else if !ep.System {
		return fmt.Errorf("current endpoint is not a system defined endpoint")
	}

	if err := ep.RemoveKeyserver(uri); err != nil {
		return err
	}

	// truncating file before writing new contents and syncing to commit file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("while truncating remote config file: %s", err)
	}

	if n, err := file.Seek(0, io.SeekStart); err != nil || n != 0 {
		return fmt.Errorf("failed to reset %s cursor: %s", file.Name(), err)
	}

	if _, err := c.WriteTo(file); err != nil {
		return fmt.Errorf("while writing remote config to file: %s", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to flush remote config file %s: %s", file.Name(), err)
	}

	return nil
}
