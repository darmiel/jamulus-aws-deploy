package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type InstallJamulusEC2Menu *EC2Menu

func NewInstallJamulusMenu(ec *ec2.EC2, instance *ec2.Instance, parent *Menu) InstallJamulusEC2Menu {
	menu := &EC2Menu{
		Menu: &Menu{Parent: parent},
		ec:   ec,
	}
	menu.Print = func() {

		// logging in with ssh
		user := "ec2-user"
		host := *instance.PublicDnsName
		port := 22

		log.Printf("Connecting to %s@%s:%d\n", user, host, port)

		// get path of key pair
		var q survey.Prompt

		log.Println("In order for us to continue setting up the instance, I need the path to the key file.")
		q = &survey.Input{
			Message: "Path of your key pair",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}

		var err error
		var signer ssh.Signer
		for {
			// TODO: remove debug path
			var keyPairPath string = "/Users/dstatzne/.ssh/mac.pem"
			if keyPairPath == "" {
				if err := survey.AskOne(q, &keyPairPath); err != nil {
					log.Fatalln("Error reading your answer:", err)
					return
				}
			}

			var data []byte
			if data, err = ioutil.ReadFile(keyPairPath); err != nil {
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
		log.Println("Read ssh key.")

		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		var client *ssh.Client
		for {
			if client, err = ssh.Dial("tcp", host+":22", config); err != nil {
				log.Println("Error connecting to", host, ". Retrying ... (", err, ")")
				time.Sleep(time.Second)
				continue
			}
			break
		}
		defer client.Close()

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
