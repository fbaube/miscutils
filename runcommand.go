package miscutils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

// RunCommand runs the command with 0-or-more arguments. It returns the
// command's Stdout and Stderr, plus any error that prevented execution.
// Note that the command and its arguments must be passed as a list of
// individual tokens, rather than as a single string with blanks as
// separators.
func RunCommand(cmdname string, args ...string) (stdout string, stderr string, err error) {
	fmt.Println("RunCommand:", cmdname, args)
	// e.g. cmd := exec.Command("ls", "-lah")
	cmd := exec.Command(cmdname, args...)
	var stout, sterr bytes.Buffer
	cmd.Stdout = &stout
	cmd.Stderr = &sterr
	err = cmd.Run()
	if err != nil {
		err = errors.Wrap(err, "RunCommand failed")
	}
	return string(stout.Bytes()), string(sterr.Bytes()), err
}
