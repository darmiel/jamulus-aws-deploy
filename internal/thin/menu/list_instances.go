package menu

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"time"
)

type Menu struct {
	ec *ec2.EC2
}

func NewMenu(ec *ec2.EC2) *Menu {
	return &Menu{ec}
}

const (
	Refresh   = "🚀️| Refresh"
	CreateNew = "🎉 | Deploy new instance"
)

func (m *Menu) DisplayListInstances() {
	// load instances
	resp, err := m.ec.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String(common.JamulusDefHeader)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: aws.StringSlice([]string{"pending", "running", "shutting-down", "stopping", "stopped"}),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// make options
	optMap := make(map[string]*ec2.Instance)
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			title := fmt.Sprintf("[%s] %s (%s) [running for %s]",
				common.GetPrettyState(i.State), *i.InstanceId, *i.PublicIpAddress, time.Since(*i.LaunchTime))
			optMap[title] = i
		}
	}

	opts := make([]string, len(optMap)+2)
	opts[0] = Refresh
	opts[1] = CreateNew
	i := 2
	for k := range optMap {
		opts[i] = k
		i++
	}

	id := common.Select("Select action", opts)
	switch id {
	case Refresh:
		m.DisplayListInstances()
		return
	case CreateNew:
		m.DisplayDeployNew()
		return
	default:
		instance := optMap[id]
		if instance == nil {
			panic("instance was empty")
		}
		m.DisplayControlInstance(instance)
		return
	}
}
