package menu

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/ctl"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/melbahja/goph"
)

const (
	ControlActionStartJamulus    = "🚀 | Start Jamulus"
	ControlActionStopJamulus     = "🔻 | Stop Jamulus"
	ControlActionToggleRecording = "🎤 | Toggle Recording"
	ControlActionTerminate       = "🗑 | Terminate"
)

func GetTemplate(instance *ec2.Instance) *templates.Template {
	// get template
	var tplName string
	for _, tag := range instance.Tags {
		if *tag.Key == common.JamulusTemplateHeader {
			tplName = *tag.Value
			break
		}
	}
	if tplName == "" {
		fmt.Println(common.ERRPrefix(), "Could not determine local template!")
		return nil
	}
	// parse template
	tpl, err := templates.FromFile(tplName)
	if err != nil {
		fmt.Println(common.ERRPrefix(), err)
		return nil
	}
	fmt.Println(common.AWSPrefix(), "Using template:", tpl.Template.Name, "by", tpl.Template.Author)
	return tpl
}

func (m *Menu) DisplayControlInstance(instance *ec2.Instance) {
	tpl := GetTemplate(instance)
	if tpl == nil {
		return
	}
	action := common.Select("Select action", []string{
		ControlActionStartJamulus,
		ControlActionStopJamulus,
		ControlActionToggleRecording,
		ControlActionTerminate,
	})

	switch action {
	case ControlActionTerminate:
		_, err := m.ec.TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{instance.InstanceId},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(common.AWSPrefix(), "Terminating", *instance.InstanceId, "...")
		return
	}

	fmt.Println(common.SSHPrefix(), "waiting for ssh ...")

	// get ssh
	key, err := goph.Key(tpl.Instance.KeyPair.LocalPath, "")
	if err != nil {
		fmt.Println(common.ERRPrefix(), "error loading key:", err)
		return
	}
	client, err := goph.NewUnknown("ec2-user", *instance.PublicIpAddress, key)
	ssh := sshc.Must(client, err)

	switch action {
	case ControlActionStartJamulus:
		ctl.StartJamulus(ssh, tpl)
		return

	case ControlActionStopJamulus:
		ctl.StopJamulus(ssh, true, tpl)
		return

	case ControlActionToggleRecording:
		containers := ssh.DockerPs(templates.JamulusDockerImage)
		if len(containers) == 0 {
			fmt.Println(common.ERRPrefix(), "there are no jamulus servers running")
			return
		}
		ctl.JamulusRecord(ssh, common.Select("select server to toggle recording", containers),
			ctl.JamulusRecModeToggle, true)
	}
}
