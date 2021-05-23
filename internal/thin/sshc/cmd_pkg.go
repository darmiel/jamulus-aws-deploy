package sshc

import "strings"

func (s *SSHC) YumUpdate() {
	s.MustExecute("sudo yum update -y")
}

func (s *SSHC) YumInstall(pkg ...string) bool {
	if len(pkg) == 0 {
		return false
	}
	return s.IsStatusCode(0, "sudo yum install %s -y", strings.Join(pkg, " "))
}

func (s *SSHC) PkgWhich(pkg string) bool {
	return s.IsStatusCode(0, "which %s", pkg)
}

const (
	ServiceStatusRunning  = 0
	ServiceStatusInactive = 3
)

func (s *SSHC) SystemCtlStatus(pkg string) int {
	return s.MustExecute("sudo systemctl status %s", pkg).StatusCode
}

func (s *SSHC) SystemCtlIsRunning(pkg string) bool {
	return s.SystemCtlStatus(pkg) == ServiceStatusRunning
}

func (s *SSHC) SystemCtl(pkg, action string) bool {
	return s.IsStatusCode(0, "sudo systemctl %s %s", action, pkg)
}
