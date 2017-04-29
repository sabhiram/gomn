package rpc

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

type RPC struct {
	pivxCLIPath string // path to the pivx-cli tool
	options     string // path to the data directory
}

func New() (*RPC, error) {
	return &RPC{
		pivxCLIPath: "/Users/shaba/Desktop/work/code/github/crypto/PIVX/src/pivx-cli",
		options:     "-datadir=/Users/shaba/.crypto/pivx",
	}, nil
}

func (r *RPC) rpcCmd(cmds ...string) ([]byte, []byte, error) {
	if len(r.options) > 0 {
		cmds = append([]string{r.options}, cmds...)
	}

	cmd := exec.Command(r.pivxCLIPath, cmds...)
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	stdPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	stderr, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return nil, nil, err
	}
	stdout, err := ioutil.ReadAll(stdPipe)
	if err != nil {
		return nil, nil, err
	}
	return stdout, stderr, err
}

func (r *RPC) GetInfo() (interface{}, error) {
	stdout, stderr, err := r.rpcCmd("getinfo")
	if err != nil {
		return nil, err
	}

	fmt.Printf("GETINFO:\n  StdOut: %s\n  StdErr: %s\n", string(stdout), string(stderr))
	return nil, nil
}

func (r *RPC) MasternodeStatus() (interface{}, error) {
	stdout, stderr, err := r.rpcCmd("masternode", "status")
	if err != nil {
		return nil, err
	}

	fmt.Printf("GETINFO:\n  StdOut: %s\n  StdErr: %s\n", string(stdout), string(stderr))
	return nil, nil
}
