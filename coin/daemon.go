package coin

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

////////////////////////////////////////////////////////////////////////////////

// ExecCmd takes a executable at path specified by `cmd`, where `args` are a
// set of options to execute the command with.  Returns the stdout, stderr and
// any errors from trying to execute the command.  This function blocks until
// the command finishes.
func ExecCmd(cmd string, args ...string) (io.ReadCloser, io.ReadCloser, error) {
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
	return stdPipe, errPipe, err
}

////////////////////////////////////////////////////////////////////////////////

func (c *Coin) StartDaemon() error {
	stdout, stderr, err := ExecCmd(c.GetDaemonBinPath())
	go func(outp, errp io.ReadCloser) {
		_, err := ioutil.ReadAll(outp)
		if err != nil {
			fmt.Printf("Unable to read stdout from cmd! %s\n", err.Error())
		}
		_, err = ioutil.ReadAll(errp)
		if err != nil {
			fmt.Printf("Unable to read stderr from cmd! %s\n", err.Error())
		}
	}(stdout, stderr)
	return err
}

////////////////////////////////////////////////////////////////////////////////
