package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/melbahja/goph"
	"github.com/muesli/termenv"
	"log"
	"strconv"
)

const (
	Region = "eu-central-1"
)

var tpl = templates.Must(templates.FromFile("InstanceTemplate.json"))

func main() {
	key, err := goph.Key(tpl.Instance.KeyPair.LocalPath, "")
	if err != nil {
		log.Fatalln("error loading key:", err)
		return
	}
	client, err := goph.NewUnknown("ec2-user", "52.28.85.232", key)
	ssh := sshc.Must(client, err)
	StartJamulus(ssh)
}

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
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "Select servers to stop",
		Options: running,
	}, &stop); err != nil {
		panic(err)
	}
	for _, id := range stop {
		common.PrintSSHState("Shutting down #" + id).
			Report(ssh.DockerContainerStop(id, 25))
	}
}

func StartJamulus(ssh *sshc.SSHC) {
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
		fmt.Println(common.JAMPrefix(), "Requesting logs for", common.Color(id, termenv.ANSIBlue.Sequence(false)))
		resp = ssh.MustExecute("sudo docker logs " + id)
		if resp.StatusCode != 0 {
			fmt.Println(common.ERRPrefix(), string(resp.Data))
		} else {
			common.LvlPrint(common.JAMPrefix(), string(resp.Data))
		}
	}
}

/*
func _main() {
	// create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		log.Fatalln("Error creating session:", err)
		return
	}
	ec := ec2.New(sess, aws.NewConfig().WithRegion(Region))

	//
	s := tsess.Session{Template: tpl, EC2: ec}

	d, _ := json.Marshal(tpl.Instance)
	log.Println("Creating instance:", string(d), "...")
	if _, err := s.CreateInstances(); err != nil {
		log.Fatalln("error creating instances:", err)
	}
}
*/
