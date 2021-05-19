package tpl

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	JamulusStatusHeader  = "Jamulus-Status"
	JamulusStatusCreated = "Created"
	JamulusStatusDone    = "Done"
)

func Is(instance *ec2.Instance, state string) bool {
	for _, tag := range instance.Tags {
		if tag.Key != nil && *tag.Key == JamulusStatusHeader {
			if tag.Value != nil && *tag.Value == state {
				return true
			}
		}
	}
	return false
}
