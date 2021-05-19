package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"log"
	"time"
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

		// security group
		var securityGroupId *string
		securityGroup := "JamulusSVR"
		groups, err := ec.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			log.Fatalln("Error reading security groups:", err)
			return
		}
		for _, sg := range groups.SecurityGroups {
			if *sg.GroupName == securityGroup {
				securityGroupId = sg.GroupId
				break
			}
		}

		s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)

		// create security group
		if securityGroupId == nil {
			s.Prefix = "🤔 Creating security group ... "
			s.FinalMSG = "😁 Created security group!"
			s.Start()

			// create security group
			cresp, err := ec.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
				Description:       aws.String("Allows SSH and the Jamulus default port"),
				GroupName:         aws.String("JamulusSVR"),
				TagSpecifications: nil,
				VpcId:             nil,
			})
			if err != nil {
				log.Fatalln("Error creating security group:", err)
				return
			}
			securityGroupId = cresp.GroupId

			// assign rules
			if _, err := ec.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: securityGroupId,
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

		s.Prefix = "🤔 Crating instance ..."
		s.FinalMSG = "😁 Created instance!"
		s.Start()
		// create instance
		runInput := &ec2.RunInstancesInput{
			ImageId:          aws.String(CreateAMI),
			InstanceType:     aws.String(instanceType),
			MinCount:         aws.Int64(1),
			MaxCount:         aws.Int64(1),
			KeyName:          aws.String(keyName),
			SecurityGroupIds: []*string{securityGroupId},
		}
		rresp, err := ec.RunInstances(runInput)
		if err != nil {
			log.Fatalln("Error creating instance:", err)
			return
		}
		instance := rresp.Instances[0]
		// attach tag
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
		s.Stop()
		fmt.Println()

		NewInstallJamulusMenu(ec, instance, menu.Menu).Print()
	}
	return menu
}
