package menu

import (
	"fmt"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/ctl"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/waiter"
)

func (m *Menu) DisplayDeployNew() {
	// select template
	tpl := templates.SelectTemplate()
	fmt.Println(common.AWSPrefix(), "Using template", tpl.Template.Name, "(", tpl.Template.Description, ")")

	sess := tsess.NewTemplatedSession(m.ec, tpl)

	//
	instances, err := sess.CreateInstances()
	if err != nil {
		panic(err)
	}
	if len(instances) != 1 {
		fmt.Println(common.ERRPrefix(), "Invalid instance count:", len(instances))
		return
	}
	instance := instances[0]
	//

	fmt.Println(common.AWSPrefix(), "Created instance:", *instance.InstanceId)
	fmt.Println(common.AWSPrefix(), "Waiting until instance is", common.Color("running", "#A8CC8C"))
	if instance, err = waiter.WaitForState(m.ec, *instance.InstanceId, "running"); err != nil {
		panic(err)
	}

	fmt.Println(common.AWSPrefix(), "Instance is running. Waiting for SSH ...")
	ssh, err := waiter.WaitForSSHInstance(instance, tpl)
	if err != nil {
		panic(err)
	}

	// start Jamulus
	ctl.StartJamulus(ssh, tpl)
}
