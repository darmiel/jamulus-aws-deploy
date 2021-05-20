package menus

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"os"
	"path"
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

		// save template
		if !temp.Template.IsTemplate {
			var saveTemplate bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Save [Server] Template?",
				Default: false,
			}, &saveTemplate); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}

			// ask for name and description
			if err := survey.Ask([]*survey.Question{
				{
					Name:     "TemplateName",
					Prompt:   &survey.Input{Message: "Template Name"},
					Validate: survey.Required,
				},
				{
					Name:   "TemplateDescription",
					Prompt: &survey.Input{Message: "Template Description"},
				},
			}, temp); err != nil {
				log.Fatalln("Error reading your answer:", err)
				return
			}

			// encode to json
			data, err := json.Marshal(temp)
			if err != nil {
				log.Fatalln("Error encoding to JSON:", err)
				return
			}

			// generate uuid
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

		NewInstallJamulusMenu(ec, instance, temp, menu.Menu).Print()
	}
	return menu
}
