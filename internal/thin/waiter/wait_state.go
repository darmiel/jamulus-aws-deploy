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
	sp := common.NewSpinner(
		fmt.Sprintf("%s 🤔 Waiting for instance to be '%s'", common.AWSPrefix(), state),
		fmt.Sprintf("%s 😁 Instance is '%s'", common.AWSPrefix(), state))
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

		// Spinner style
		sp.Prefix = fmt.Sprintf("%s 🤔 Waiting for instance to be '%s' [%s] ",
			common.AWSPrefix(), state, common.GetPrettyState(i.State))

		if *i.State.Name != state {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	sp.Stop()
	fmt.Println()
	return
}
