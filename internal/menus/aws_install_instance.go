package menus

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"os"
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

		client, err := temp.OpenSession(ec, instance)
		if err != nil {
			log.Fatalln("Error opening client:", err)
			return
		}
		defer client.Close()

		sess, err := client.NewSession()
		if err != nil {
			log.Fatalln("Error opening socket:", err)
			return
		}
		defer sess.Close()

		const command = "sudo yum update -y; sudo yum install docker -y; sudo service docker start; sudo docker run --rm hello-world"

		log.Println("Running command: `" + command + "`")

		sess.Stdout = os.Stdout
		sess.Stderr = os.Stderr

		if err := sess.Run(command); err != nil {
			log.Fatalln("error running command:", err)
			return
		}

		log.Println("ok!")
	}
	return menu
}
