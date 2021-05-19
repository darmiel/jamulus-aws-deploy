package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"os"
	"path/filepath"
	"time"
)

type CreateInstanceEC2Menu *EC2Menu

const (
	CreateAMI = "ami-043097594a7df80ec"

	DefaultKeyPair       = "jamulus-cert"
	DefaultSecurityGroup = "jamulus-security-group"
)

func NewCreateInstanceMenu(ec *ec2.EC2, parent *Menu) CreateInstanceEC2Menu {
	menu := &EC2Menu{
		ec:   ec,
		Menu: &Menu{Parent: parent},
	}
	menu.Print = func() {
		temp := tpl.SelectTemplate()
		if temp == nil {
			return
		}

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

		// INSTANCE TYPE
		if temp.InstanceType == "" {
			q = &survey.Select{
				Message: "Select instance type",
				Options: []string{
					ec2.InstanceTypeT2Micro,
					ec2.InstanceTypeC5Large,
					ec2.InstanceTypeC5Xlarge,
					ec2.InstanceTypeC52xlarge,
				},
			}
			if err := survey.AskOne(q, &temp.InstanceType); err != nil {
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
			if temp.KeyPair == "" {
				temp.KeyPair = DefaultKeyPair
			}
			// ask to create key pair?
			pair, err := ec.CreateKeyPair(&ec2.CreateKeyPairInput{
				KeyName: aws.String(temp.KeyPair),
			})
			if err != nil {
				log.Fatalln("Error creating key pair:", err)
				return
			}

			// save key
			var outPath string
			if temp.KeyPairPath != "" {
				if _, err := os.Stat(temp.KeyPairPath); os.IsNotExist(err) {
					outPath = temp.KeyPairPath
				}
			}

			// ask for path
			if outPath == "" {
				q = &survey.Input{
					Message: "Path of your key pair",
					Suggest: func(toComplete string) []string {
						files, _ := filepath.Glob(toComplete + "*")
						return files
					},
				}
				if err := survey.AskOne(q, &outPath); err != nil {
					log.Fatalln("Error reading your answer:", err)
					return
				}
			}

			// save to path
			if err := os.WriteFile(outPath, []byte(*pair.KeyMaterial), 0644); err != nil {
				log.Fatalln("Error writing cert:", err)
				return
			}
		} else {
			if temp.KeyPair == "" {
				opts := make([]string, len(resp.KeyPairs))
				i := 0
				for _, pair := range resp.KeyPairs {
					opts[i] = *pair.KeyName
					i++
				}
				q = &survey.Select{Message: "Select Key-Pair", Options: opts}
				if err := survey.AskOne(q, &temp.KeyPair); err != nil {
					log.Fatalln("Error reading key pair from input:", err)
					return
				}
			}
		}

		// security group
		groups, err := ec.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			log.Fatalln("Error reading security groups:", err)
			return
		}
		for _, sg := range groups.SecurityGroups {
			if *sg.GroupName == DefaultSecurityGroup {
				temp.SecurityGroupID = *sg.GroupId
				break
			}
		}

		s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)

		// create security group
		if temp.SecurityGroupID == "" {
			s.Prefix = "🤔 Creating security group "
			s.FinalMSG = "😁 Created security group!"
			s.Start()

			// create security group
			createResponse, err := ec.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
				Description:       aws.String("Allows SSH and the Jamulus default port"),
				GroupName:         aws.String(DefaultSecurityGroup),
				TagSpecifications: nil,
				VpcId:             nil,
			})
			if err != nil {
				log.Fatalln("Error creating security group:", err)
				return
			}
			temp.SecurityGroupID = *createResponse.GroupId

			// assign rules to security group
			if _, err := ec.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: &temp.SecurityGroupID,
				IpPermissions: []*ec2.IpPermission{
					(&ec2.IpPermission{}).
						SetIpProtocol("tcp").
						SetFromPort(22).
						SetToPort(22).
						SetIpRanges([]*ec2.IpRange{
							(&ec2.IpRange{}).
								SetCidrIp("0.0.0.0/0"),
						}),
					(&ec2.IpPermission{}).
						SetIpProtocol("udp").
						SetFromPort(22124).
						SetToPort(22124).
						SetIpRanges([]*ec2.IpRange{
							(&ec2.IpRange{}).
								SetCidrIp("0.0.0.0/0"),
						}),
				},
			}); err != nil {
				log.Fatalln("Error assigning rules to security group:", err)
				return
			}
			s.Stop()
			fmt.Println()
		}

		s.Prefix = "🤔 Crating instance "
		s.FinalMSG = "😁 Created instance!"
		s.Start()

		// create instance
		runInput := &ec2.RunInstancesInput{
			ImageId:          aws.String(CreateAMI),
			InstanceType:     aws.String(temp.InstanceType),
			MinCount:         aws.Int64(1),
			MaxCount:         aws.Int64(1),
			KeyName:          aws.String(temp.KeyPair),
			SecurityGroupIds: []*string{&temp.SecurityGroupID},
		}
		rresp, err := ec.RunInstances(runInput)
		if err != nil {
			log.Fatalln("Error creating instance:", err)
			return
		}
		instance := rresp.Instances[0]

		// attach tags
		tagInput := &ec2.CreateTagsInput{
			Resources: aws.StringSlice([]string{*instance.InstanceId}),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Jamulus"),
					Value: aws.String("Yes"),
				},
				{
					Key:   aws.String(tpl.JamulusStatusHeader),
					Value: aws.String(tpl.JamulusStatusCreated),
				},
			},
		}
		if _, err := ec.CreateTags(tagInput); err != nil {
			log.Fatalln("Error attaching tags:", err)
		}

		s.Stop()
		fmt.Println()

		NewInstallJamulusMenu(ec, instance, temp, menu.Menu).Print()
	}
	return menu
}
