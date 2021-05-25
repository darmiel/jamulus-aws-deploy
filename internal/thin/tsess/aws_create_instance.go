package tsess

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"log"
)

const (
	Owner = "Unknown" // change with `go build -ldflags "-X tsess.Owner=<name>"`
)

func (s *Session) CreateInstances() (instances []*ec2.Instance, err error) {
	var sg *ec2.SecurityGroup
	log.Println("find security group ...")
	for {
		if sg, err = s.FindSecurityGroup(); err != nil {
			if err == ErrSecurityGroupNotFound {
				log.Println("creating security group")
				err = s.CreateSecurityGroup()
				fmt.Println(err)
				continue
			}
			return
		}
		break
	}

	log.Println("find key pair ...")
	var kp *ec2.KeyPairInfo
	for {
		if kp, err = s.FindKeyPair(); err != nil {
			if err == ErrKeyPairNotFound {
				log.Println("creating and saving key pair")
				err = s.CreateKeyPair()
				fmt.Println(err)
				continue
			}
			return
		}
		break
	}
	fmt.Println("found key pair:", *kp.KeyPairId)

	var resv *ec2.Reservation
	if resv, err = s.EC2.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(s.Instance.AMI),
		InstanceType:     aws.String(s.Instance.Type),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		KeyName:          aws.String(*kp.KeyName),
		SecurityGroupIds: []*string{sg.GroupId},
	}); err != nil {
		return
	}

	tags := []*ec2.Tag{
		&common.JamulusDefTag,
		{
			Key:   aws.String(common.JamulusOwnerHeader),
			Value: aws.String(Owner),
		},
	}

	// append local template
	if c := s.Template.LocalTemplate; c != "" {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(common.JamulusTemplateHeader),
			Value: aws.String(c),
		})
	}

	if err = s.AttachTags(resv.Instances, tags); err != nil {
		return
	}
	return resv.Instances, nil
}

func (s *Session) AttachTags(instances []*ec2.Instance, tags []*ec2.Tag) (err error) {
	resources := make([]*string, len(instances))
	for idx, i := range instances {
		resources[idx] = i.InstanceId
	}
	_, err = s.EC2.CreateTags(&ec2.CreateTagsInput{
		Resources: resources,
		Tags:      tags,
	})
	return
}
