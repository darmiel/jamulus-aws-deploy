package tsess

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
)

type Session struct {
	*templates.Template
	*ec2.EC2
}

func NewTemplatedSession(ec *ec2.EC2, tpl *templates.Template) *Session {
	return &Session{
		Template: tpl,
		EC2:      ec,
	}
}
