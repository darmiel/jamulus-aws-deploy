package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"time"
)

const (
	Region = "eu-central-1"
)

func GetPrettyState(state *ec2.InstanceState) string {
	switch *state.Name {
	case "pending":
		return "⏱ pending"
	case "running":
		return "✅ running"
	case "shutting-down":
		return "🔻 shutting down"
	case "terminated":
		return "🗑 terminated"
	case "stopping":
		return "🥱 stopping"
	case "stopped":
		return "😴 stopped"
	}
	return ""
}

func main() {
	// create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		log.Fatalln("Error creating session:", err)
		return
	}

	ec := ec2.New(sess, aws.NewConfig().WithRegion(Region))
	PrintMenu(ec)
}

func PrintMenu(ec *ec2.EC2) {
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
		fmt.Println("*", *r.ReservationId)
		for _, i := range r.Instances {
			opt[fmt.Sprintf("[%s] %s [running for %s]", GetPrettyState(i.State), *i.InstanceId, time.Now().Sub(*i.LaunchTime).String())] = i
		}
	}

	opts := make([]string, len(opt)+1)
	i := 0
	opts[i] = "♻️ Refresh"
	for id, _ := range opt {
		i++
		opts[i] = id
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

	// (re-) print
	if instanceId == "♻️ Refresh" {
		PrintMenu(ec)
		return
	}

	instance, ok := opt[instanceId]
	if !ok {
		log.Fatalln("Invalid instance selected")
		return
	}

	// what to do?
	q = &survey.Select{
		Message: "Select Action",
		Options: []string{"👋 Go Back", "Terminate Instance", "Stop Instance", "Start Instance"},
	}

	var action string
	if err := survey.AskOne(q, &action); err != nil {
		log.Fatalln("Error selecting action:", err)
		return
	}

	switch action {
	case "Terminate Instance":
		log.Println("Terminating", *instance.InstanceId, "...")
		resp, err := ec.TerminateInstances(&ec2.TerminateInstancesInput{
			DryRun:      aws.Bool(false),
			InstanceIds: []*string{instance.InstanceId},
		})
		if err != nil {
			log.Fatalln("Error terminating instance:", err)
			return
		}
		log.Println("Done! (", resp.String(), ")")
		break

	case "Stop Instance":
		log.Println("Stopping", *instance.InstanceId, "...")
		resp, err := ec.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: aws.StringSlice([]string{*instance.InstanceId}),
		})
		if err != nil {
			log.Fatalln("Error stopping instance:", err)
			return
		}
		log.Println("Done! (", resp.String(), ")")
		break

	case "Start Instance":
		log.Println("Starting", *instance.InstanceId, "...")
		resp, err := ec.StartInstances(&ec2.StartInstancesInput{
			InstanceIds: aws.StringSlice([]string{*instance.InstanceId}),
		})
		if err != nil {
			log.Fatalln("Error starting instance:", err)
			return
		}
		log.Println("Done! (", resp.String(), ")")
		break
	case "👋 Go Back":
		break
	}
	PrintMenu(ec)
}
