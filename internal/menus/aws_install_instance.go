package menus

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"github.com/melbahja/goph"
	"log"
	"time"
)

type InstallJamulusEC2Menu *EC2Menu

func NewInstallJamulusMenu(ec *ec2.EC2, instance *ec2.Instance, temp *tpl.CreateInstanceTemplate, parent *Menu) InstallJamulusEC2Menu {
	menu := &EC2Menu{
		Menu: &Menu{Parent: parent},
		ec:   ec,
	}
	menu.Print = func() {
		if temp == nil {
			if temp = tpl.SelectTemplate(tpl.TemplateTypeInstance); temp == nil {
				return
			}
		}

		host := temp.WaitForHost(ec, instance)
		if host == "" {
			log.Fatalln("empty host")
			return
		}

		key, err := goph.Key(temp.Instance.KeyPairPath, "")
		if err != nil {
			log.Fatalln("error loading key:", err)
			return
		}

		var ssh *sshc.SSHC

		// wait until instance is running
		s := tpl.NewSpinner("🤔 Waiting for SSH", "😁 SSH connected")
		for {
			client, err := goph.NewUnknown("ec2-user", host, key)
			if err != nil {
				s.Prefix = "🤔 Waiting for SSH (" + err.Error() + ") "
				time.Sleep(2 * time.Second)
				continue
			}
			ssh = sshc.Must(client, err)
			break
		}
		s.Stop()
		fmt.Println()

		// ---

		// make update
		s = tpl.NewSpinner("> 🎖 Yum Update", "> 🎖 Yum Update (Done)")
		if err := ssh.YumUpdate(); err != nil {
			log.Fatalln("  > error updating yum:", err)
			return
		}
		s.Stop()
		fmt.Println()

		// check if docker is installed
		if ok, _ := ssh.IsInstalled("docker"); !ok {
			// install docker
			s = tpl.NewSpinner("> 🎖 Installing Docker", "> 🎖 Installing Docker (Done)")
			if err := ssh.YumInstall("docker"); err != nil {
				log.Fatalln("  > error installing docker:", err)
				return
			}
			s.Stop()
			fmt.Println()
		}

		// check if docker is running
		if ok, _ := ssh.ServiceRunning("docker"); !ok {
			// start docker service
			s = tpl.NewSpinner("> 🎖 Starting Docker Service", "> 🎖 Starting Docker Service (Done)")
			if err := ssh.ServiceCtl("docker", "start"); err != nil {
				log.Fatalln("  > error starting docker service:", err)
				return
			}
			s.Stop()
			fmt.Println()
		}

		log.Println("Server set up.")

		// TODO: start Jamulus setup
		NewAskJamulusParamsMenu(menu.Menu).Print()
	}
	return menu
}
