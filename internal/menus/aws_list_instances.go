package menus

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"time"
)

const (
	Refresh   = "🚀️ | Refresh"
	CreateNew = "🎉 | Deploy new instance"
	GoBack    = "👋 | Go Back"
)

///

type ListInstancesEC2Menu *EC2Menu

func NewListInstancesEC2Menu(ec *ec2.EC2) ListInstancesEC2Menu {
	menu := &EC2Menu{
		ec:   ec,
		Menu: &Menu{},
	}
	menu.Print = func() {
		// get instances
		instances, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag-key"),
					Values: []*string{aws.String("Jamulus")},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: aws.StringSlice([]string{"pending", "running", "shutting-down", "stopping", "stopped"}),
				},
			},
		})
		if err != nil {
			log.Fatalln(err)
			return
		}

		opt := make(map[string]*ec2.Instance)
		for _, r := range instances.Reservations {
			for _, i := range r.Instances {
				isDone := ""
				if tpl.Is(i, tpl.JamulusStatusCreated) {
					isDone = "🐣 "
				}

				opt[fmt.Sprintf("%s[%s] %s (%s) [running for %s]",
					isDone,
					tpl.GetPrettyState(i.State),
					*i.InstanceId,
					*i.PublicIpAddress,
					time.Since(*i.LaunchTime).String())] = i
			}
		}

		opts := make([]string, len(opt)+2)
		opts[0] = Refresh
		opts[1] = CreateNew
		i := 2
		for id, _ := range opt {
			opts[i] = id
			i++
		}

		// ask what to do
		q := &survey.Select{
			Message: "Select Instance",
			Options: opts,
		}

		var instanceId string
		if err := survey.AskOne(q, &instanceId); err != nil {
			log.Fatalln("Error selecting instance:", err)
			return
		}

		switch instanceId {
		case Refresh:
			menu.Print()
			return

		case CreateNew:
			NewCreateInstanceMenu(ec, menu.Menu).Print()
			return

		default:
			instance, ok := opt[instanceId]
			if !ok {
				log.Fatalln("Invalid instance selected")
				return
			}
			// print control menu
			NewControlInstanceMenu(ec, instance, menu.Menu).Print()
			return
		}
	}
	return menu
}
