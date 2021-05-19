package menus

import (
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
		q := &survey.Confirm{Message: "Create (another) instance?", Default: false}
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

		runInput := &ec2.RunInstancesInput{
			ImageId:      aws.String(CreateAMI),
			InstanceType: aws.String(ec2.InstanceTypeT2Micro),
			MinCount:     aws.Int64(1),
			MaxCount:     aws.Int64(1),
		}

		resp, err := ec.RunInstances(runInput)
		if err != nil {
			log.Fatalln("Error creating instance:", err)
			return
		}

		// attach tag
		tagInput := &ec2.CreateTagsInput{
			Resources: aws.StringSlice([]string{*resp.Instances[0].InstanceId}),
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
		menu.Print()
	}
	return menu
}
