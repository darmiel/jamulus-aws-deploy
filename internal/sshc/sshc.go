package sshc

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
)

type SSHC struct {
	*goph.Client
}

func Must(g *goph.Client, err error) *SSHC {
	if err != nil {
		log.Fatalln("error connecting to ssh:", err)
		return nil
	}
	return &SSHC{g}
}

func (s *SSHC) IsInstalled(prog string) (bool, error) {
	_, err := s.Run(fmt.Sprintf("which %s", prog))
	if err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			if ee.ExitStatus() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (s *SSHC) YumUpdate() (err error) {
	_, err = s.Run("sudo yum update -y")
	return
}

func (s *SSHC) YumInstall(pkg string) (err error) {
	_, err = s.Run(fmt.Sprintf("sudo yum install %s -y", pkg))
	return
}

func (s *SSHC) ServiceCtl(service, action string) (err error) {
	_, err = s.Run(fmt.Sprintf("sudo service %s %s", service, action))
	return
}

func (s *SSHC) ServiceRunning(service string) (resp bool, err error) {
	_, err = s.Run(fmt.Sprintf("sudo service %s status", service))
	if err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			if ee.ExitStatus() > 0 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (s *SSHC) DockerContainerExists(id string) (resp bool, err error) {
	_, err = s.Run(fmt.Sprintf("sudo docker inspect %s", id))
	if err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			if ee.ExitStatus() > 0 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
