package ctl

import (
	"fmt"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/muesli/termenv"
	"strings"
)

func StartJamulus(ssh *sshc.SSHC, tpl *templates.Template) {
	// check if docker is installed
	if !ssh.PkgWhich("docker") {
		// sudo yum update -y
		common.PrintSSHState("Yum Update").
			Report(ssh.YumUpdate())

		// sudo yum install docker -y
		common.PrintSSHState("Install Docker").
			Report(ssh.YumInstall("docker"))

		// sudo systemctl start docker
		common.PrintSSHState("Start Docker Service").
			Report(ssh.SystemCtl("docker", "start"))
	}

	// mkdirs
	if c := tpl.Jamulus.LogPath; c != "" {
		if !ssh.DirExists(c) {
			common.PrintSSHState("mkdir (log) " + common.Color(c, "#A8CC8C").String()).
				Report(ssh.DirCreate(c))
		}
	}
	if c := tpl.Jamulus.Recording.Path; c != "" {
		if !ssh.DirExists(c) {
			common.PrintSSHState("mkdir (rec) " + common.Color(c, "#A8CC8C").String()).
				Report(ssh.DirCreate(c))
		}
	}

	// stop old jamulus servers
	StopJamulus(ssh, false)

	// start
	cmd := tpl.Jamulus.CreateArgs()
	fmt.Println(common.SSHPrefix(), common.Color(cmd, "#D290E4"))

	// sudo docker run [ ... ]
	state := common.PrintSSHState("Starting Jamulus Server")
	resp := ssh.MustExecute(cmd)
	state.Report(resp)

	// sudo docker logs [ ... ]
	if resp.StatusCode == 0 {
		id := string(resp.Data)
		if strings.Contains(id, "Pull") {
			split := strings.Split(id, "\n")
			id = strings.TrimSpace(split[len(split)-2])
		}

		fmt.Println(common.JAMPrefix(), "Requesting logs for", common.Color(id, termenv.ANSIBlue.Sequence(false)))
		resp = ssh.MustExecute("sudo docker logs " + id)
		if resp.StatusCode != 0 {
			fmt.Println(common.ERRPrefix(), string(resp.Data))
		} else {
			common.LvlPrint(common.JAMPrefix(), string(resp.Data))
		}
	}
	fmt.Println(common.JAMPrefix(), "Connect Jamulus to",
		common.Color(ssh.Client().Config.Addr+":22124", "#A8CC8C"))
}
