package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"path/filepath"
)

type CreateInstanceEC2Menu *EC2Menu

func NewCreateInstanceMenu(ec *ec2.EC2, parent *Menu) CreateInstanceEC2Menu {
	menu := &EC2Menu{
		ec:   ec,
		Menu: &Menu{Parent: parent},
	}
	menu.Print = func() {
		temp := tpl.SelectTemplate(tpl.TemplateTypeInstance)
		if temp == nil {
			return
		}

		var q survey.Prompt
		q = &survey.Confirm{
			Message: "Create (another) instance?",
			Default: true,
		}
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

		// INSTANCE TYPE
		if temp.Instance.InstanceType == "" {
			q = &survey.Select{
				Message: "Select instance type",
				Options: []string{
					ec2.InstanceTypeT2Micro,
					ec2.InstanceTypeC5Large,
					ec2.InstanceTypeC5Xlarge,
					ec2.InstanceTypeC52xlarge,
				},
			}
			if err := survey.AskOne(q, &temp.Instance.InstanceType); err != nil {
				log.Fatalln("Error reading instance type:", err)
				return
			}
		}

		// KEYPAIR NAME
		resp, err := ec.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
		if err != nil {
			log.Fatalln("Error reading your key pairs from aws:", err)
			return
		}
		if len(resp.KeyPairs) == 0 {
			if temp.Instance.KeyPair == "" {
				temp.Instance.KeyPair = tpl.DefaultKeyPair
			}

			// create key pair
			temp.CreateKeyPair(ec)
		} else if temp.Instance.KeyPair == "" {
			opts := make([]string, len(resp.KeyPairs))
			i := 0
			for _, pair := range resp.KeyPairs {
				opts[i] = *pair.KeyName
				i++
			}
			q = &survey.Select{Message: "Select Key-Pair", Options: opts}
			if err := survey.AskOne(q, &temp.Instance.KeyPair); err != nil {
				log.Fatalln("Error reading key pair from input:", err)
				return
			}
		}

		// SECURITY GROUP
		groups, err := ec.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			log.Fatalln("Error reading security groups:", err)
			return
		}
		for _, sg := range groups.SecurityGroups {
			if *sg.GroupName == tpl.DefaultSecurityGroup {
				temp.Instance.SecurityGroupID = *sg.GroupId
				break
			}
		}
		// create?
		if temp.Instance.SecurityGroupID == "" {
			temp.CreateSecurityGroup(ec)
			fmt.Println()
		}

		instance := temp.CreateInstance(ec)
		fmt.Println()

		// get key pair
		if temp.Instance.KeyPairPath == "" {
			if err := survey.AskOne(&survey.Input{
				Message: "Path of your key pair",
				Suggest: func(toComplete string) []string {
					files, _ := filepath.Glob(toComplete + "*")
					return files
				},
			}, &temp.Instance.KeyPairPath); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}
		}

		// ask to save template
		temp.AskSave()

		NewInstallJamulusMenu(ec, instance, temp, menu.Menu).Print()
	}
	return menu
}
