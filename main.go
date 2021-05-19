package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/menus"
	"log"
)

const (
	Region = "eu-central-1"
)

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

	menu := menus.NewListInstancesEC2Menu(ec)
	menu.Print()
}
