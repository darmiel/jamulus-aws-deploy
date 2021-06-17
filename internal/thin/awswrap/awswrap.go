package awswrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
)

type AWSWrap struct {
	CredPath string
	Region   string
	Config   *aws.Config
	STS      *sts.STS
	EC2      *ec2.EC2
}

func GetInstanceOwner(instance *ec2.Instance) string {
	for _, t := range instance.Tags {
		if t == nil || t.Key == nil || t.Value == nil {
			continue
		}
		if *t.Key == common.JamulusOwnerHeader {
			return *t.Value
		}
	}
	return ""
}
