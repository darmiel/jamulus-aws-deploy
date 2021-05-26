package ctl

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"strconv"
)

func StopJamulus(ssh *sshc.SSHC, verbose bool) {
	running := ssh.DockerPs(templates.JamulusDockerImage)
	if len(running) <= 0 {
		if verbose {
			fmt.Println(common.SSHPrefix(), "No Jamulus servers running")
		}
		return
	}
	fmt.Println(common.SSHPrefix(), "There are/is",
		common.Color(strconv.Itoa(len(running)), "#E88388"),
		"Jamulus server/s running")

	stop := make([]string, 0)

	// ask which one to stop
	if len(running) == 1 {
		var stopit bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Stop Jamulus server #" + running[0],
			Default: true,
		}, &stopit); err != nil {
			panic(err)
		}
		if !stopit {
			return
		}
		stop = append(stop, running[0])
	} else {
		if err := survey.AskOne(&survey.MultiSelect{
			Message: "Select servers to stop",
			Options: running,
		}, &stop); err != nil {
			panic(err)
		}
	}

	// stop
	for _, id := range stop {
		common.PrintSSHState("Shutting down #" + id).
			Report(ssh.DockerContainerStop(id, 25))
	}
}
