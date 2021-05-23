package sshc

func (s *SSHC) DockerImageExists(image string) bool {
	// TODO: Implement
	return false
}

func (s *SSHC) DockerContainerStatus(container string) string {
	// TODO: Implement
	return ""
}

func (s *SSHC) DockerContainerRunning(container string) bool {
	// TODO: Implement
	return false
}

func (s *SSHC) DockerSendSignal(container, signal string) {
	// TODO: Implement
}

func (s *SSHC) DockerContainerExec(container, command string) {
	// TODO: Implement
}

func (s *SSHC) DockerContainerLogs(container string) string {
	// TODO: Implement
	return ""
}

func (s *SSHC) DockerContainerStart() {
	// TODO: Implement
}

func (s *SSHC) DockerContainerStop() {
	// TODO: Implement
}

func (s *SSHC) DockerContainerRestart(container string) bool {
	// TODO: Implement
	return false
}
