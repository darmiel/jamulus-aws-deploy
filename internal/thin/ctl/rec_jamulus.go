package ctl

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
)

const (
	JamulusRecModeToggle = iota
	JamulusRecModeStart
)

func JamulusRecord(ssh *sshc.SSHC, container string, mode int, verbose bool) {
	var signal string
	switch mode {
	case JamulusRecModeStart:
		signal = "SIGUSR1"
	case JamulusRecModeToggle:
		signal = "SIGUSR2"
	default:
		if verbose {
			fmt.Println(common.ERRPrefix(), "invalid mode.")
		}
		return
	}

	if container == "" {
		running := ssh.DockerPs(templates.JamulusDockerImage)
		if len(running) == 0 {
			if verbose {
				fmt.Println(common.ERRPrefix(), "There is no Jamulus server running")
			}
			return
		} else if len(running) == 1 {
			container = running[0]
		} else {
			q := &survey.Select{
				Message: "Select server to send " + signal,
				Options: running,
			}
			if err := survey.AskOne(q, &container); err != nil {
				panic(err)
			}
		}
	}
	fmt.Print(common.SSHPrefix(), " Sending signal ",
		common.Color(signal, "#DBAB79"), " to ", common.Color(container, "#66C2CD"),
		" ... ")
	if resp := ssh.DockerSendSignal(container, signal); resp.StatusCode == 0 {
		fmt.Println("👍")
	} else {
		fmt.Println("🤬")
		common.LvlPrint(common.ERRPrefix(), string(resp.Data))
	}
}
