package sshc

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	DockerStatusRunning    = "running"
	DockerStatusPaused     = "paused"
	DockerStatusRestarting = "restarting"
	DockerStatusDead       = "dead"
	DockerStatusUnknown    = "unknown"
)

var (
	dockerContainerIdRegex = regexp.MustCompile("(?m)^([a-f0-9]+)")
)

func (s *SSHC) DockerContainerStatus(container string) string {
	resp := s.MustExecute("sudo docker inspect -f '{{ .State.Status }}' %s", strconv.Quote(container))
	if resp.StatusCode != 0 {
		return DockerStatusUnknown
	}
	return string(resp.Data)
}

func (s *SSHC) DockerPs(image string) []string {
	resp := s.MustExecute("sudo docker ps -a | grep %s", strconv.Quote(image))
	res := make([]string, 0)
	if resp.StatusCode != 0 {
		return res
	}
	return dockerContainerIdRegex.FindAllString(string(resp.Data), -1)
}

func (s *SSHC) DockerContainerRunning(container string) bool {
	return s.DockerContainerStatus(container) == DockerStatusRunning
}

func (s *SSHC) DockerSendSignal(container, signal string) *SSHCCommandResult {
	return s.MustExecute("sudo docker kill -s %s %s", strconv.Quote(signal), strconv.Quote(container))
}

func (s *SSHC) DockerContainerExec(container, command string, it bool) *SSHCCommandResult {
	args := make([]string, 0)
	args = append(args, "sudo docker exec")
	if it {
		args = append(args, "-it")
	}
	args = append(args, strconv.Quote(container))
	args = append(args, strconv.Quote(command))
	return s.MustExecute(strings.Join(args, " "))
}

func (s *SSHC) DockerContainerLogs(container string) string {
	resp := s.MustExecute("sudo docker logs %s", strconv.Quote(container))
	if resp.StatusCode != 0 {
		return ""
	}
	return string(resp.Data)
}

func (s *SSHC) DockerContainerStop(container string, timeout int) *SSHCCommandResult {
	return s.MustExecute("sudo docker stop -t %d %s", timeout, strconv.Quote(container))
}

func (s *SSHC) DockerContainerRestart(container string, timeout int) *SSHCCommandResult {
	return s.MustExecute("sudo docker container restart -t %d %s", timeout, strconv.Quote(container))
}
