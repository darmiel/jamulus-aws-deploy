package menus

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

const (
	StartInstance     = "🔌 | Start"
	StopInstance      = "🔻 | Shut Down"
	TerminateInstance = "🗑 | Terminate"
)

type ControlInstanceEC2Menu *EC2Menu

func NewControlInstanceMenu(ec *ec2.EC2, instance *ec2.Instance, parent *Menu) ControlInstanceEC2Menu {
	menu := &EC2Menu{
		ec: ec,
		Menu: &Menu{
			Parent: parent,
		},
	}

	menu.Print = func() {

		// build options
		options := make([]string, 1)
		options[0] = GoBack
		state := *instance.State.Name
		if state != ec2.InstanceStateNameTerminated {
			// start instance
			if state != ec2.InstanceStateNameRunning {
				options = append(options, StartInstance)
			}
			// stop instance
			if state != ec2.InstanceStateNameStopping && state != ec2.InstanceStateNameStopped {
				options = append(options, StopInstance)
			}
			// terminate instance
			options = append(options, TerminateInstance)
		}

		// what to do?
		var action string
		q := &survey.Select{Message: "Select Action", Options: options}
		if err := survey.AskOne(q, &action); err != nil {
			log.Fatalln("Error selecting action:", err)
			return
		}

		switch action {
		case TerminateInstance:
			log.Println("Terminating", *instance.InstanceId, "...")
			_, err := ec.TerminateInstances(&ec2.TerminateInstancesInput{
				DryRun:      aws.Bool(false),
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				log.Fatalln("Error terminating instance:", err)
				return
			}
			log.Println("Done!")
			break

		case StopInstance:
			log.Println("Stopping", *instance.InstanceId, "...")
			_, err := ec.StopInstances(&ec2.StopInstancesInput{
				InstanceIds: aws.StringSlice([]string{*instance.InstanceId}),
			})
			if err != nil {
				log.Fatalln("Error stopping instance:", err)
				return
			}
			log.Println("Done!")
			break

		case StartInstance:
			log.Println("Starting", *instance.InstanceId, "...")
			_, err := ec.StartInstances(&ec2.StartInstancesInput{
				InstanceIds: aws.StringSlice([]string{*instance.InstanceId}),
			})
			if err != nil {
				log.Fatalln("Error starting instance:", err)
				return
			}
			log.Println("Done!")
			break
		case GoBack:
			break
		}

		menu.Back()
	}

	return menu
}
