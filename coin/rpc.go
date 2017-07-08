package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"io/ioutil"
	"os/exec"
)

////////////////////////////////////////////////////////////////////////////////

// RPCCommand takes a executable at path specified by `cmd`, where `args` are a
// set of options to execute the command with.  Returns the stdout, stderr and
// any errors from trying to execute the command.  This function blocks until
// the command finishes.
func RPCCommand(cmd string, args ...string) ([]byte, []byte, error) {
	c := exec.Command(cmd, args...)
	stdPipe, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	errPipe, err := c.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := c.Start(); err != nil {
		return nil, nil, err
	}

	stdout, err := ioutil.ReadAll(stdPipe)
	if err != nil {
		return nil, nil, err
	}
	stderr, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return nil, nil, err
	}
	return stdout, stderr, err
}

////////////////////////////////////////////////////////////////////////////////
