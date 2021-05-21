package menus

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"strconv"
	"strings"
)

type AskJamulusParamsMenu *Menu

func NewAskJamulusParamsMenu(parent *Menu, ssh *sshc.SSHC, ec *ec2.EC2, instance *ec2.Instance) AskJamulusParamsMenu {
	menu := &Menu{
		Parent: parent,
	}
	menu.Print = func() {
		temp := tpl.SelectTemplate(tpl.TemplateTypeJamulus)
		if temp == nil {
			log.Fatalln("error selecting template")
			return
		}

		// make public?
		if temp.Jamulus.CentralServer == nil {
			if common.Bool("Make Server Public", false) {
				// central server
				temp.Jamulus.CentralServer = aws.String(common.FlatSelect("Select Server", map[string]string{
					"Any Genre 1":             "anygenre1.jamulus.io:22124",
					"Any Genre 2":             "anygenre2.jamulus.io:22224",
					"Any Genre 3":             "anygenre3.jamulus.io:22624",
					"Genre Rock":              "rock.jamulus.io:22424",
					"Genre Jazz":              "jazz.jamulus.io:22324",
					"Genre Classical/Folk":    "classical.jamulus.io:22524",
					"Genre Choral/Barbershop": "choral.jamulus.io:22724",
				}))

				// server info
				data := struct {
					Name   string
					City   string
					Locale string
				}{}
				q := []*survey.Question{
					{
						Name:     "Name",
						Prompt:   &survey.Input{Message: "Public Server-Name"},
						Validate: survey.Required,
					},
					{
						Name:     "City",
						Prompt:   &survey.Input{Message: "Public Server-City"},
						Validate: survey.Required,
					},
					{
						Name:     "Locale",
						Prompt:   &survey.Input{Message: "Public Server-Locale"},
						Validate: survey.Required,
					},
				}

				if err := survey.Ask(q, &data); err != nil {
					log.Fatalln("error reading your answer:", err)
					return
				}

				temp.Jamulus.ServerInfo = aws.String(fmt.Sprintf("%s;%s;%s", data.Name, data.City, data.Locale))
			} else {
				temp.Jamulus.CentralServer = aws.String("-")
			}
		}

		// max users
		if temp.Jamulus.MaxUsers == nil {
			usersStr := common.StringValidate("Max Users", "25", func(ans interface{}) error {
				s, ok := ans.(string)
				if !ok {
					return errors.New("not a string")
				}
				_, err := strconv.Atoi(s)
				return err
			})
			atoi, err := strconv.Atoi(usersStr)
			if err != nil {
				log.Fatalln("error reading max users:", err)
				return
			}
			temp.Jamulus.MaxUsers = &atoi
		}

		// welcome message
		if temp.Jamulus.WelcomeMessage == nil {
			if temp.Jamulus.WelcomeMessage = aws.String(common.String("Welcome Message", "Hallo!")); *temp.Jamulus.WelcomeMessage == "" {
				temp.Jamulus.WelcomeMessage = aws.String("-")
			}
		}

		// fast update
		if temp.Jamulus.FastUpdate == nil {
			temp.Jamulus.FastUpdate = aws.Bool(common.Bool("Use Fast Update?", true))
		}

		// multi threading
		if temp.Jamulus.EnableMultiThreading == nil {
			temp.Jamulus.EnableMultiThreading = aws.Bool(common.Bool("Enable Multi Threading", true))
		}

		// recording
		if temp.Jamulus.RecordingPath == nil {
			if common.Bool("Enable Recording", true) {
				temp.Jamulus.RecordingPath = aws.String(common.String("Recording Path", "/tmp/recordings"))
				temp.Jamulus.NoRecordOnStart = aws.Bool(common.Bool("No Record On Start", true))
			} else {
				temp.Jamulus.RecordingPath = aws.String("-")
				temp.Jamulus.NoRecordOnStart = aws.Bool(true)
			}
		}

		// logging
		if temp.Jamulus.LogPath == nil {
			if common.Bool("Enable Logging", true) {
				temp.Jamulus.LogPath = aws.String(common.String("Log Path", "/tmp/logs"))
			} else {
				temp.Jamulus.LogPath = aws.String("-")
			}
		}

		temp.AskSave()

		/////////////
		// check if docker process is running
		if resp, _ := ssh.DockerContainerExists("jamulus"); resp {
			fmt.Println("> 🎖 Stopping old docker instance ...")
			_, _ = ssh.Run("sudo docker stop jamulus")
		}

		// build cli
		var b strings.Builder
		b.WriteString("sudo docker run -d --name jamulus --rm grundic/jamulus")

		// params
		// central server
		if temp.Jamulus.CentralServer != nil && *temp.Jamulus.CentralServer != "-" {
			b.WriteString(" --centralserver ")
			b.WriteString(*temp.Jamulus.CentralServer)
		}

		// server info
		if temp.Jamulus.ServerInfo != nil && *temp.Jamulus.ServerInfo != "-" {
			b.WriteString(" --serverinfo ")
			b.WriteRune('"')
			b.WriteString(*temp.Jamulus.ServerInfo)
			b.WriteRune('"')
		}

		// fast update
		if temp.Jamulus.FastUpdate != nil && *temp.Jamulus.FastUpdate {
			b.WriteString(" --fastupdate")
		}

		// log path
		if temp.Jamulus.LogPath != nil && *temp.Jamulus.LogPath != "-" {
			b.WriteString(" --log ")
			b.WriteRune('"')
			b.WriteString(*temp.Jamulus.LogPath)
			b.WriteRune('"')

			// mkdir
			var cmd = "mkdir -p \"" + *temp.Jamulus.LogPath + "\""
			fmt.Println("> 🎖", cmd)
			_, _ = ssh.Run(cmd)
		}

		// recording path
		if temp.Jamulus.RecordingPath != nil && *temp.Jamulus.RecordingPath != "-" {
			b.WriteString(" --recording ")
			b.WriteRune('"')
			b.WriteString(*temp.Jamulus.RecordingPath)
			b.WriteRune('"')

			// mkdir
			var cmd = "mkdir -p \"" + *temp.Jamulus.RecordingPath + "\""
			fmt.Println("> 🎖", cmd)
			_, _ = ssh.Run(cmd)

			if temp.Jamulus.NoRecordOnStart != nil && *temp.Jamulus.NoRecordOnStart {
				b.WriteString(" --norecord")
			}
		}

		// enable multi threading
		if temp.Jamulus.EnableMultiThreading != nil && *temp.Jamulus.EnableMultiThreading {
			b.WriteString(" --multithreading")
		}

		// max users
		if temp.Jamulus.MaxUsers != nil {
			b.WriteString(" --numchannels ")
			b.WriteString(strconv.Itoa(*temp.Jamulus.MaxUsers))
		}

		// welcome
		if temp.Jamulus.WelcomeMessage != nil && *temp.Jamulus.WelcomeMessage != "-" {
			b.WriteString(" --welcomemessage ")
			b.WriteRune('"')
			b.WriteString(*temp.Jamulus.WelcomeMessage)
			b.WriteRune('"')
		}

		fmt.Println("> 🎖 Executing:", b.String(), "...")

		resp, err := ssh.Run(b.String())
		if err != nil {
			log.Fatalln("error starting server:", err)
			return
		}
		fmt.Println("✅ Container started:", string(resp))

		// update tags
		// attach tags
		tagInput := &ec2.CreateTagsInput{
			Resources: aws.StringSlice([]string{*instance.InstanceId}),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(tpl.JamulusStatusHeader),
					Value: aws.String(tpl.JamulusStatusDone),
				},
			},
		}
		if _, err := ec.CreateTags(tagInput); err != nil {
			log.Fatalln("Error attaching tags:", err)
		}

		var host string
		if instance == nil || instance.PublicIpAddress == nil {
			host = "unknown"
		} else {
			host = *instance.PublicIpAddress
		}
		fmt.Println("🤗 Connect to", host+":22124")
	}
	return menu
}
