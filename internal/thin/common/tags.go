package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// headers
const (
	JamulusDefHeader      = "Jamulus"
	JamulusOwnerHeader    = "Jamulus-Owner"
	JamulusTemplateHeader = "Jamulus-Local-Template"
)

var (
	JamulusDefTag = ec2.Tag{
		Key:   aws.String(JamulusDefHeader),
		Value: aws.String("Yes"),
	}
)
