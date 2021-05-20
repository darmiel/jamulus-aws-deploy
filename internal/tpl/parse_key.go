package tpl

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"time"
)

func GetPrettyState(state *ec2.InstanceState) string {
	switch *state.Name {
	case "pending":
		return "⏱ | pending"
	case "running":
		return "✅ | running"
	case "shutting-down":
		return "🔻 | shutting down"
	case "terminated":
		return "🗑 | terminated"
	case "stopping":
		return "🥱 | stopping"
	case "stopped":
		return "😴 | stopped"
	}
	return ""
}

func (t *CreateInstanceTemplate) WaitForHost(ec *ec2.EC2, instance *ec2.Instance) string {
	var instanceHost string

	// wait until instance is running
	s := NewSpinner("🤔 Waiting for instance to be ready", "😁 Instance is running!")
	for {
		resp, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{instance.InstanceId},
		})
		if err != nil {
			log.Fatalln("Error reading instance:", err)
			return ""
		}
		i := resp.Reservations[0].Instances[0]
		if *i.State.Name != ec2.InstanceStateNameRunning {
			time.Sleep(2 * time.Second)
			continue
		}
		instanceHost = *i.PublicDnsName
		s.Prefix = fmt.Sprintf("🤔 Waiting for instance to be ready [%s] ", GetPrettyState(i.State))
		break
	}
	s.Stop()
	fmt.Println()

	return instanceHost
}
