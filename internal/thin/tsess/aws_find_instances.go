package tsess

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
)

func (s *Session) FindInstances(owner string) (resp []*ec2.Instance, err error) {
	filters := []*ec2.Filter{
		{
			Name:   aws.String("tag-key"),
			Values: []*string{aws.String("Jamulus")},
		},
		{
			Name:   aws.String("instance-state-name"),
			Values: aws.StringSlice([]string{"pending", "running", "shutting-down", "stopping", "stopped"}),
		},
	}
	if owner != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:" + common.JamulusOwnerHeader),
			Values: aws.StringSlice([]string{owner}),
		})
	}
	var out *ec2.DescribeInstancesOutput
	if out, err = s.EC2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: filters,
	}); err != nil {
		return
	}
	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			resp = append(resp, i)
		}
	}
	return
}
