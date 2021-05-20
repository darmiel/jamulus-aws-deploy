package tpl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"log"
	"time"
)

const CreateAMI = "ami-043097594a7df80ec"

func (t *CreateInstanceTemplate) CreateInstance(ec *ec2.EC2) *ec2.Instance {
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Prefix = "🤔 Crating instance "
	s.FinalMSG = "😁 Created instance!"
	s.Start()
	defer s.Stop()

	// create instance
	runInput := &ec2.RunInstancesInput{
		ImageId:          aws.String(CreateAMI),
		InstanceType:     aws.String(t.Instance.InstanceType),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String(t.Instance.KeyPair),
		SecurityGroupIds: []*string{&t.Instance.SecurityGroupID},
	}
	resp, err := ec.RunInstances(runInput)
	if err != nil {
		log.Fatalln("Error creating instance:", err)
		return nil
	}

	instance := resp.Instances[0]

	// attach tags
	tagInput := &ec2.CreateTagsInput{
		Resources: aws.StringSlice([]string{*instance.InstanceId}),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Jamulus"),
				Value: aws.String("Yes"),
			},
			{
				Key:   aws.String(JamulusStatusHeader),
				Value: aws.String(JamulusStatusCreated),
			},
		},
	}
	if _, err := ec.CreateTags(tagInput); err != nil {
		log.Fatalln("Error attaching tags:", err)
	}
	return instance
}
