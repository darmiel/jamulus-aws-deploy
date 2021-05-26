package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/menu"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"log"
)

const (
	Region = "eu-central-1"
)

var tpl = templates.Must(templates.FromFile("InstanceTemplate.json"))

/*
func main() {

	StartJamulus(ssh)
}
*/

func main() {
	fmt.Println("aws-deploy - compiled for", tsess.Owner)

	// create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		log.Fatalln("Error creating session:", err)
		return
	}
	ec := ec2.New(sess, aws.NewConfig().WithRegion(Region))

	m := menu.NewMenu(ec)
	m.DisplayListInstances()
}
