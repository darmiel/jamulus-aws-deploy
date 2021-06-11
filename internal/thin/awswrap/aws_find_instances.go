package awswrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
)

func (wr *AWSWrap) FindInstances(owner string, viewAll, checkJamulus bool) (resp []*ec2.Instance, err error) {
	filters := make([]*ec2.Filter, 1)

	// filter instance state
	filters[0] = &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: aws.StringSlice([]string{"pending", "running", "shutting-down", "stopping", "stopped"}),
	}

	if !viewAll {
		if checkJamulus {
			filters = append(filters, &ec2.Filter{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String("Jamulus")},
			})
		}
		if owner != "" {
			filters = append(filters, &ec2.Filter{
				Name:   aws.String("tag:" + common.JamulusOwnerHeader),
				Values: aws.StringSlice([]string{owner}),
			})
		}
	}

	var out *ec2.DescribeInstancesOutput
	if out, err = wr.EC2.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: filters,
	}); err != nil {
		return
	}
	for _, r := range out.Reservations {
		resp = append(resp, r.Instances...)
	}
	return
}
