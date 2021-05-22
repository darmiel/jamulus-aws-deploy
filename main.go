package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"log"
	"path"
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

	//

	tplPath := path.Join("flat-tpl", "InstanceTemplate.json")
	tpl, err := templates.FromFile(tplPath)
	if err != nil {
		log.Fatalln("Error reading template:", err)
		return
	}

	s := tsess.Session{Template: tpl, EC2: ec, TemplatePath: tplPath}

	d, _ := json.Marshal(tpl.Instance)
	log.Println("Creating instance:", string(d), "...")
	if _, err := s.CreateInstances(); err != nil {
		log.Fatalln("error creating instances:", err)
	}
}
