package waiter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/melbahja/goph"
	"time"
)

func WaitForSSHInstance(instance *ec2.Instance, tpl *templates.Template) (c *sshc.SSHC, err error) {
	key, err := goph.Key(tpl.Instance.KeyPair.LocalPath, "")
	if err != nil {
		return nil, err
	}
	return WaitForSSH("ec2-user", *instance.PublicIpAddress, key)
}

func WaitForSSH(user, addr string, auth goph.Auth) (c *sshc.SSHC, err error) {
	// wait until instance is running
	prefix := common.SSHPrefix().String() + " 🤔 Waiting for SSH to be ready"
	sp := common.NewSpinner(prefix,
		common.SSHPrefix().String()+" 😁 SSH available!")
	try := 1

	var client *goph.Client
	for {
		if client, err = goph.NewUnknown(user, addr, auth); err != nil {
			// TODO: Check error
			sp.Prefix = fmt.Sprintf("%s | (t=%d) wait err(%T): %v", prefix, try, err, err)
			try++
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	sp.Stop()
	fmt.Println()

	c = sshc.Must(client, err)
	return
}
