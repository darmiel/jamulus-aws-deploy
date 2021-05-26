package ctl

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"strconv"
)

const (
	JamulusRecModeToggle = iota
	JamulusRecModeStart
)

func JamulusRecord(ssh *sshc.SSHC, mode int, verbose bool) {
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

	toggleSvr := make([]string, 0)

	// ask which one to toggleSvr
	if len(running) == 1 {
		var recordit bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Toggle recording on Jamulus server #" + running[0],
			Default: true,
		}, &recordit); err != nil {
			panic(err)
		}
		if !recordit {
			return
		}
		toggleSvr = append(toggleSvr, running[0])
	} else {
		if err := survey.AskOne(&survey.MultiSelect{
			Message: "Select servers to toggle recording",
			Options: running,
		}, &toggleSvr); err != nil {
			panic(err)
		}
	}

	for _, container := range toggleSvr {
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
}
