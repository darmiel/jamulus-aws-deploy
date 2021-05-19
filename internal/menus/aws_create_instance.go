package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

type CreateInstanceEC2Menu *EC2Menu

const (
	CreateAMI = "ami-043097594a7df80ec"
)

func NewCreateInstanceMenu(ec *ec2.EC2, parent *Menu) CreateInstanceEC2Menu {
	menu := &EC2Menu{
		ec:   ec,
		Menu: &Menu{Parent: parent},
	}
	menu.Print = func() {
		var q survey.Prompt
		q = &survey.Confirm{Message: "Create (another) instance?", Default: false}
		var createNew bool
		if err := survey.AskOne(q, &createNew); err != nil {
			log.Fatalln("Error selecting:", err)
			return
		}
		// fallback
		if !createNew {
			menu.Back()
			return
		}

		// TODO: Implement template

		// ask for image
		var instanceType string
		q = &survey.Select{
			Message: "Select instance type",
			Options: []string{
				ec2.InstanceTypeT2Micro,
				ec2.InstanceTypeC5Large,
				ec2.InstanceTypeC5Xlarge,
				ec2.InstanceTypeC52xlarge,
			},
		}
		if err := survey.AskOne(q, &instanceType); err != nil {
			log.Fatalln("Error reading intance type:", err)
			return
		}

		// ask for key name
		var keyName string
		// read key names
		resp, err := ec.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
		if err != nil {
			log.Fatalln("Error reading your key pairs from aws:", err)
			return
		}
		if len(resp.KeyPairs) == 0 {
			keyName = ""

			log.Println("WARNING :: You don't have any key pairs on your AWS account")
			log.Println("WARNING :: If you don't link a key pair to the server,")
			log.Println("WARNING :: we won't be able to connect!")
			log.Println("WARNING :: Do you really want to create an instance anyways?")

			fmt.Println()
			log.Println("WARNING :: You can create a key pair here:")
			log.Println("WARNING :: https://console.aws.amazon.com/ec2/v2/home#KeyPairs:")
			fmt.Println()

			q = &survey.Confirm{Message: "Create Instance Anyways?", Default: false}
			var createAnyways bool
			if err := survey.AskOne(q, &createAnyways); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}
			if !createAnyways {
				menu.Back()
				return
			}
			log.Println("Okay. Good luck!")
		} else {
			opts := make([]string, len(resp.KeyPairs))
			i := 0
			for _, pair := range resp.KeyPairs {
				opts[i] = *pair.KeyName
				i++
			}
			q = &survey.Select{Message: "Select Key-Pair", Options: opts}
			if err := survey.AskOne(q, &keyName); err != nil {
				log.Fatalln("Error reading key pair from input:", err)
				return
			}
		}

		// TODO: security group

		// create instance
		log.Println("Creating instance ...")
		runInput := &ec2.RunInstancesInput{
			ImageId:          aws.String(CreateAMI),
			InstanceType:     aws.String(instanceType),
			MinCount:         aws.Int64(1),
			MaxCount:         aws.Int64(1),
			KeyName:          aws.String(keyName),
			SecurityGroupIds: []*string{aws.String("sg-807df1e1")}, // default
		}

		rresp, err := ec.RunInstances(runInput)
		if err != nil {
			log.Fatalln("Error creating instance:", err)
			return
		}
		instance := rresp.Instances[0]

		// attach tag
		log.Println("Attaching tag to", *instance.InstanceId, "...")
		tagInput := &ec2.CreateTagsInput{
			Resources: aws.StringSlice([]string{*instance.InstanceId}),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Jamulus"),
					Value: aws.String("Yes"),
				},
			},
		}
		if _, err := ec.CreateTags(tagInput); err != nil {
			log.Fatalln("Error attaching tags:", err)
		}

		log.Println("✅ Created!")

		NewInstallJamulusMenu(ec, instance, menu.Menu).Print()
	}
	return menu
}
