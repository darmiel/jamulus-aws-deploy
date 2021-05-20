package menus

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"github.com/melbahja/goph"
	"log"
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

		ssh, err := goph.NewUnknown("ec2-user", host, key)
		if err != nil {
			log.Fatalln("error connecting:", err)
			return
		}

		// TODO: commands here
		log.Println(ssh.Run("whiami"))

		if err := ssh.Close(); err != nil {
			log.Fatalln("error closing connection: ", err)
		}

		log.Println("ok!")
	}
	return menu
}
