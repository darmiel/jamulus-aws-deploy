package menu

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/awswrap"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"time"
)

type Menu struct {
	ec *ec2.EC2
	wr *awswrap.AWSWrap
}

func NewMenu(ec *ec2.EC2, wr *awswrap.AWSWrap) *Menu {
	return &Menu{ec, wr}
}

const (
	Refresh      = "🚀 | Refresh"
	CreateNew    = "🎉 | Deploy new instance"
	ShowMoreLess = "🔍 | Show more/less"
)

func (m *Menu) DisplayListInstances(owner string, showAll, checkJamulus bool) {
	fmt.Println(common.AWSPrefix(), "Loading instances ...")

	// load instances
	resp, err := m.wr.FindInstances(owner, showAll, checkJamulus)
	if err != nil {
		panic(err)
	}

	// make options
	optMap := make(map[string]*ec2.Instance)
	for _, i := range resp {
		title := fmt.Sprintf("💻 | [%s] %s (%s) [running for %s]",
			common.GetPrettyState(i.State), *i.InstanceId, *i.PublicIpAddress, time.Since(*i.LaunchTime))
		optMap[title] = i
	}

	opts := make([]string, len(optMap)+3)
	opts[0] = Refresh
	opts[1] = ShowMoreLess
	opts[2] = CreateNew
	i := 2
	for k := range optMap {
		opts[i] = k
		i++
	}

	// Select Deploy on start
	// if an instance was found, select instance
	var def interface{} = CreateNew
	for t := range optMap {
		def = t
		break
	}
	id := common.Select("Select action", opts, def)

	switch id {
	case Refresh:
		break

	case ShowMoreLess:

		opts = []string{
			"🐣 | Show only my own Jamulus instances",
			"🐣 | Show all my own instances",
			"🏘 | Show all Jamulus instances",
			"🏘 | Show all instances",
		}

		q := &survey.Select{
			Message: "Toggle filters",
			Options: opts,
			Default: opts[0],
		}

		if err = survey.AskOne(q, &id); err != nil {
			panic(err)
		}

		switch id {
		case opts[0]:
			owner = common.Owner
			checkJamulus = true
			showAll = false
		case opts[1]:
			owner = common.Owner
			checkJamulus = false
			showAll = false
		case opts[2]:
			owner = ""
			checkJamulus = true
			showAll = false
		case opts[3]:
			showAll = true
		}

	case CreateNew:
		m.DisplayDeployNew()

	default:
		instance := optMap[id]
		if instance == nil {
			panic("instance was empty")
		}
		m.DisplayControlInstance(instance)
	}

	m.DisplayListInstances(owner, showAll, checkJamulus)
}
