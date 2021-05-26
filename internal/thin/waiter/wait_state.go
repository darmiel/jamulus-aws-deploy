package waiter

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"time"
)

func WaitForState(ec *ec2.EC2, instanceId, state string) (i *ec2.Instance, err error) {
	// wait until instance is running
	sp := common.NewSpinner("🤔 Waiting for instance to be ready", "😁 Instance is running!")
	for {
		var resp *ec2.DescribeInstancesOutput
		if resp, err = ec.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instanceId}),
		}); err != nil {
			return
		}
		for _, resv := range resp.Reservations {
			for _, inst := range resv.Instances {
				i = inst
				break
			}
		}
		if i == nil {
			return nil, errors.New("no instance returned")
		}
		if *i.State.Name != ec2.InstanceStateNameRunning {
			time.Sleep(2 * time.Second)
			continue
		}
		sp.Prefix = fmt.Sprintf("🤔 Waiting for instance to be ready [%s] ", common.GetPrettyState(i.State))
		break
	}
	sp.Stop()
	fmt.Println()
	return
}
