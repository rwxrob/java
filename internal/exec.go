// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Exec checks for existence of first argument as an executable on the
// system and then runs it without exiting in a way that is supported
// across all architectures that Go supports. The stdin, stdout, and stderr are
// connected directly to that of the calling program. Sometimes this is
// insufficient and the UNIX-specific SysExec is preferred. For example,
// when handing over control to a terminal editor such as Vim.
func Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Out returns the standard output of the executed command as
// a string. Errors are logged but not returned.
func Out(args ...string) string {
	if len(args) == 0 {
		log.Println("missing name of executable")
		return ""
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Println(err)
		return ""
	}
	out, err := exec.Command(path, args[1:]...).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
