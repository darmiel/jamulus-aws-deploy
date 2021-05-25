package tsess

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
)

type Session struct {
	*templates.Template
	*ec2.EC2
}
