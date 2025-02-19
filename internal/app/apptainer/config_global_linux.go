// Copyright (c) 2021 Apptainer a Series of LF Projects LLC
//   For website terms of use, trademark policy, privacy policy and other
//   project policies see https://lfprojects.org/policies
// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package apptainer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/apptainer/apptainer/pkg/util/apptainerconf"
	"golang.org/x/sys/unix"
)

// GlobalConfigOp defines a type for a global configuration operation.
type GlobalConfigOp uint8

const (
	// GlobalConfigSet is the operation to set a configuration directive value.
	GlobalConfigSet GlobalConfigOp = iota
	// GlobalConfigUnset is the operation to unset a configuration directive value.
	GlobalConfigUnset
	// GlobalConfigGet is the operation to get a configuration directive value.
	GlobalConfigGet
	// GlobalConfigReset is the operation to reset a configuration directive value.
	GlobalConfigReset
)

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func generateConfig(path string, directives apptainerconf.Directives, dry bool) error {
	// Generate the config structure from our directives
	c, err := apptainerconf.GetConfig(directives)
	if err != nil {
		return fmt.Errorf("configuration directive invalid: %w", err)
	}

	// Write a config file to our in memory buffer
	newConfig := new(bytes.Buffer)
	if err := apptainerconf.Generate(newConfig, "", c); err != nil {
		return fmt.Errorf("while generating configuration from template: %w", err)
	}

	// Dry run = write to Stdout
	out := os.Stdout
	// Not dry run = create / overwrite existing file, now we know we have valid content
	if !dry {
		unix.Umask(0)

		flags := os.O_CREATE | os.O_TRUNC | unix.O_NOFOLLOW | os.O_RDWR
		nf, err := os.OpenFile(path, flags, 0o644)
		if err != nil {
			return fmt.Errorf("while creating configuration file %s: %w", path, err)
		}
		defer nf.Close()
		out = nf
	}

	_, err = io.Copy(out, newConfig)
	if err != nil {
		return fmt.Errorf("while writing configuration file %s: %w", path, err)
	}

	return nil
}

// GlobalConfig allows to set/unset/reset a configuration directive value
// in apptainer.conf
func GlobalConfig(args []string, configFile string, dry bool, op GlobalConfigOp) error {
	directive := args[0]
	value := ""

	if directive == "" {
		return fmt.Errorf("you must specify a configuration directive")
	}
	if len(args) > 1 {
		value = args[1]
	}

	if !apptainerconf.HasDirective(directive) {
		return fmt.Errorf("%q is not a valid configuration directive", directive)
	}

	f, err := os.OpenFile(configFile, os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("while opening configuration file %s: %s", configFile, err)
	}
	defer f.Close()

	directives, err := apptainerconf.GetDirectives(f)
	if err != nil {
		return err
	}

	values := []string{}
	if value != "" {
		for _, v := range strings.Split(value, ",") {
			va := strings.TrimSpace(v)
			if va != "" {
				if contains(values, va) {
					continue
				}
				values = append(values, va)
			}
		}
	}

	switch op {
	case GlobalConfigSet:
		if len(values) == 0 {
			return fmt.Errorf("you must specify a value for directive %q", directive)
		}

		if _, ok := directives[directive]; ok {
			for i := len(values) - 1; i >= 0; i-- {
				if contains(directives[directive], values[i]) {
					values = append(values[:i], values[i+1:]...)
				}
			}
			directives[directive] = append(values, directives[directive]...)
		} else {
			directives[directive] = values
		}
	case GlobalConfigUnset:
		unset := false
		for i := len(directives[directive]) - 1; i >= 0; i-- {
			for _, v := range values {
				if directives[directive][i] == v {
					unset = true
					directives[directive] = append(directives[directive][:i], directives[directive][i+1:]...)
				}
			}
		}
		if !unset {
			return fmt.Errorf("value '%s' not found for directive %q", value, directive)
		}
	case GlobalConfigGet:
		if len(directives[directive]) > 0 {
			fmt.Println(strings.Join(directives[directive], ","))
		}
		return nil
	case GlobalConfigReset:
		delete(directives, directive)
	}

	return generateConfig(configFile, directives, dry)
}
