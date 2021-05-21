package tsess

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
)

func (s *Session) CreateInstances() (instances []*ec2.Instance, err error) {
	var sg *ec2.SecurityGroup
	if sg, err = s.FindSecurityGroup(); err != nil {
		return
	}
	var kp *ec2.KeyPairInfo
	if kp, err = s.FindKeyPair(); err != nil {
		return
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
	if err = s.AttachTags(resv.Instances, []*ec2.Tag{
		&common.JamulusDefTag,
		&common.JamulusStatusCreatedTag,
	}); err != nil {
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
