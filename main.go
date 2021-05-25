package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"github.com/melbahja/goph"
	"github.com/muesli/termenv"
	"log"
	"path"
	"strconv"
	"strings"
)

const (
	Region = "eu-central-1"
)

var tpl = templates.Must(templates.FromFile(path.Join("flat-tpl", "InstanceTemplate.json")))
var p = termenv.ColorProfile()

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

func sshPrefix() termenv.Style {
	return termenv.String(" SSH ").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

func errPrefix() termenv.Style {
	return termenv.String(" ERR ").Foreground(p.Color("0")).Background(p.Color("#E88388"))
}

type PrintReport struct{}

func (*PrintReport) Report(resp interface{}) {
	var st string
	switch t := resp.(type) {
	case *sshc.SSHCCommandResult:
		if t.StatusCode == 0 {
			st = "👍"
		} else {
			st = "🤬"
		}
		fmt.Println(st)
		if t.StatusCode != 0 {
			for _, line := range strings.Split(string(t.Data), "\n") {
				if len(strings.TrimSpace(line)) <= 0 {
					continue
				}
				fmt.Println(errPrefix(), termenv.String(line).Foreground(p.Color("#DBAB79")))
			}
		}
	case bool:
		if t {
			st = "👍"
		} else {
			st = "🤬"
		}
		fmt.Println(st)
	default:
		fmt.Printf("Report :: invalid type: %T (%v)\n", t, t)
	}
}

func PrintOkState(msg string) *PrintReport {
	fmt.Print(sshPrefix(), " 🔨 | ", msg, " ... ")
	return &PrintReport{}
}

func StopJamulus(ssh *sshc.SSHC, verbose bool) {
	running := ssh.DockerPs(templates.JamulusDockerImage)
	if len(running) <= 0 {
		if verbose {
			fmt.Println(sshPrefix(), "No Jamulus servers running")
		}
		return
	}
	fmt.Println(sshPrefix(), "There are/is",
		termenv.String(strconv.Itoa(len(running))).Foreground(p.Color("#E88388")),
		"Jamulus server/s running")

	stop := make([]string, 0)
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "Select servers to stop",
		Options: running,
	}, &stop); err != nil {
		panic(err)
	}
	for _, id := range stop {
		PrintOkState("Shutting down #" + id).
			Report(ssh.DockerContainerStop(id, 25))
	}
}

func StartJamulus(ssh *sshc.SSHC) {
	// check if docker is installed
	if !ssh.PkgWhich("docker") {
		// sudo yum update -y
		PrintOkState("Yum Update").
			Report(ssh.YumUpdate())

		// sudo yum install docker -y
		PrintOkState("Install Docker").
			Report(ssh.YumInstall("docker"))

		// sudo systemctl start docker
		PrintOkState("Start Docker Service").
			Report(ssh.SystemCtl("docker", "start"))
	}

	// mkdirs
	if c := tpl.Jamulus.LogPath; c != "" {
		if !ssh.DirExists(c) {
			PrintOkState("mkdir (log) " + termenv.String(c).Foreground(p.Color("#A8CC8C")).String()).
				Report(ssh.DirCreate(c))
		}
	}
	if c := tpl.Jamulus.Recording.Path; c != "" {
		if !ssh.DirExists(c) {
			PrintOkState("mkdir (rec) " + termenv.String(c).Foreground(p.Color("#A8CC8C")).String()).
				Report(ssh.DirCreate(c))
		}
	}

	// stop old jamulus servers
	StopJamulus(ssh, false)

	// start
	cmd := tpl.Jamulus.CreateArgs()
	fmt.Println(sshPrefix(), termenv.String(cmd).Foreground(p.Color("#D290E4")))

	// sudo docker run [ ... ]
	PrintOkState("Starting Jamulus Server").
		Report(ssh.MustExecute(cmd))
}

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
