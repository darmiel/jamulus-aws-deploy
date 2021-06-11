package awswrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AWSWrap struct {
	CredPath string
	Region   string
	Config   *aws.Config
	STS      *sts.STS
	EC2      *ec2.EC2
}
