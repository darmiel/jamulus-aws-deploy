package menus

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type InstallJamulusEC2Menu *EC2Menu

func NewInstallJamulusMenu(ec *ec2.EC2, instance *ec2.Instance, temp *tpl.CreateInstanceTemplate, parent *Menu) InstallJamulusEC2Menu {
	menu := &EC2Menu{
		Menu: &Menu{Parent: parent},
		ec:   ec,
	}
	menu.Print = func() {
		var q survey.Prompt
		var err error

		// get path of key pair
		q = &survey.Input{
			Message: "Path of your key pair",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}

		var signer ssh.Signer
		for {
			if temp.KeyPairPath == "" {
				if err := survey.AskOne(q, &temp.KeyPairPath); err != nil {
					log.Fatalln("Error reading your answer:", err)
					return
				}
			}

			var data []byte
			if data, err = ioutil.ReadFile(temp.KeyPairPath); err != nil {
				fmt.Println()
				log.Println("Error reading file:", err)
				log.Println("Please try again!")
				fmt.Println()
				continue
			}

			// parse key
			if signer, err = ssh.ParsePrivateKey(data); err != nil {
				fmt.Println()
				log.Println("Error parsing key:", err)
				log.Println("Please try again with another key!")
				fmt.Println()
				continue
			}
			break
		}

		// save template
		if !temp.IsTemplate {
			q = &survey.Confirm{
				Message: "Save Template?",
				Default: false,
			}
			var saveTemplate bool
			if err := survey.AskOne(q, &saveTemplate); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}
			// ask for name and description
			qu := []*survey.Question{
				{
					Name:     "TemplateName",
					Prompt:   &survey.Input{Message: "Template Name"},
					Validate: survey.Required,
				},
				{
					Name:   "TemplateDescription",
					Prompt: &survey.Input{Message: "Template Description"},
				},
			}
			if err := survey.Ask(qu, temp); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}
			// encode
			data, err := json.Marshal(temp)
			if err != nil {
				log.Fatalln("Error encoding to JSON:", err)
				return
			}

			b := make([]byte, 16)
			if _, err = rand.Read(b); err != nil {
				log.Fatalln("Error generating id:", err)
				return
			}

			uuid := fmt.Sprintf("%x-%x-%x-%x-%x.json", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

			// write
			if err := os.WriteFile(path.Join("templates", uuid), data, 0755); err != nil {
				log.Fatalln("Error writing to file:", err)
				return
			}
		}

		config := &ssh.ClientConfig{
			User: "ec2-user",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// wait until instance is running
		s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
		s.Prefix = "🤔 Waiting for instance to be ready "
		s.FinalMSG = "😁 Instance is running!"
		s.Start()

		var host string
		for {
			resp, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				log.Fatalln("Error reading instance:", err)
				return
			}
			i := resp.Reservations[0].Instances[0]
			s.Prefix = fmt.Sprintf("🤔 Waiting for instance to be ready [%s] ", GetPrettyState(i.State))
			if *i.State.Name == ec2.InstanceStateNameRunning {
				host = *i.PublicDnsName
				break
			}
			time.Sleep(2 * time.Second)
		}

		s.Stop()
		fmt.Println()

		// empty host?
		if host == "" {
			log.Fatalln("Error reading server hostname! (empty)")
			return
		}

		s.Prefix = "🤔 Waiting for SSH connection "
		s.FinalMSG = "😁 Connected to SSH"
		s.Start()

		var client *ssh.Client
		for {
			if client, err = ssh.Dial("tcp", host+":22", config); err != nil {
				s.Prefix = "🤨 Waiting for SSH connection [" + err.Error() + "] "
				time.Sleep(time.Second)
				continue
			}
			break
		}
		defer client.Close()

		s.Stop()
		fmt.Println()

		// start session
		sess, err := client.NewSession()
		if err != nil {
			log.Fatalln("Error starting session:", err)
			return
		}

		sess.Stdout = os.Stdout
		log.Println("running ls -larth")
		if err := sess.Run("ls -larth"); err != nil {
			log.Fatalln("error:", err)
		}
	}
	return menu
}
