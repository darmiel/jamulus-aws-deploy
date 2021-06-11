package tsess

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"strconv"
)

func (s *Session) CreateInstances() (instances []*ec2.Instance, err error) {
	var sg *ec2.SecurityGroup

	fmt.Println(common.AWSPrefix(), "🔍 Find security group",
		common.Color(s.Template.Instance.SecurityGroupName, "#66C2CD"))

	for {
		if sg, err = s.FindSecurityGroup(); err != nil {
			if err == ErrSecurityGroupNotFound {
				fmt.Println(common.AWSPrefix(), "✏️ Creating new security group")
				err = s.CreateSecurityGroup()
				fmt.Println(err)
				continue
			}
			return
		}
		break
	}

	fmt.Println(common.AWSPrefix(), "🔍 Find key pair",
		common.Color(s.Template.Instance.KeyPair.Name, "#66C2CD"))

	var kp *ec2.KeyPairInfo
	for {
		if kp, err = s.FindKeyPair(); err != nil {
			if err == ErrKeyPairNotFound {
				fmt.Println(common.AWSPrefix(), "✏️ Creating and saving new key pair")
				err = s.CreateKeyPair()
				fmt.Println(err)
				continue
			}
			return
		}
		break
	}

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
			Value: aws.String(common.Owner),
		},
	}

	if c := s.Template.LocalTemplate; c != "" {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(common.JamulusTemplateHeader),
			Value: aws.String(c),
		})
	}

	fmt.Println(common.AWSPrefix(), "Attaching", common.Color(strconv.Itoa(len(tags)), "#E88388"), "tags")
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
