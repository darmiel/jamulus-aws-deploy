package sshc

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type SSHCCommandResult struct {
	Data       []byte
	StatusCode int
	ExitError  string
}

func (r *SSHCCommandResult) String() string {
	if r.Data == nil {
		return ""
	} else {
		return string(r.Data)
	}
}

// yum install %s -y
func (s *SSHC) Execute(cmd string, args ...interface{}) (res *SSHCCommandResult, err error) {
	// format command
	var command string
	if args != nil && len(args) > 0 {
		command = fmt.Sprintf(cmd, args...)
	} else {
		command = cmd
	}

	// run command
	resp, runerr := s.client.Run(command)
	res = &SSHCCommandResult{Data: resp}

	if runerr != nil {
		if ee, ok := runerr.(*ssh.ExitError); ok {
			res.StatusCode = ee.ExitStatus()
			res.ExitError = ee.Error()
		} else {
			err = runerr
		}
	}

	return
}

func (s *SSHC) MustExecute(cmd string, args ...interface{}) *SSHCCommandResult {
	resp, err := s.Execute(cmd, args...)
	if err != nil {
		panic(err)
	}
	return resp
}

func (s *SSHC) IsStatusCode(status int, cmd string, args ...interface{}) bool {
	return s.MustExecute(cmd, args...).StatusCode == status
}
