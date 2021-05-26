package menu

import (
	"fmt"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
)

func (m *Menu) DisplayDeployNew() {
	// select template
	tpl := templates.SelectTemplate()
	fmt.Println(common.AWSPrefix(), "Using template", tpl.Template.Name, "(", tpl.Template.Description, ")")

	sess := tsess.Session{
		Template: tpl,
		EC2:      m.ec,
	}

	instances, err := sess.CreateInstances()
	if err != nil {
		panic(err)
	}

	for _, i := range instances {
		fmt.Println(common.AWSPrefix(), "Created instance:", *i.InstanceId)
	}
}
