package menu

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/ctl"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/waiter"
)

const (
	HeaderActionJamulus          = "Jamulus"
	ControlActionStartJamulus    = "🚀 | Start Jamulus"
	ControlActionStopJamulus     = "🔻 | Stop Jamulus"
	ControlActionToggleRecording = "🎤 | Toggle Recording"
	HeaderActionSCP              = "SCP"
	ControlActionGetRecordings   = "🎙 Browse Recordings"
	ControlActionGetLogs         = "📝 | Browse Logs"
	HeaderActionAWS              = "AWS"
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
		HeaderActionJamulus,
		ControlActionStartJamulus,
		ControlActionStopJamulus,
		ControlActionToggleRecording,

		HeaderActionSCP,
		ControlActionGetRecordings,
		ControlActionGetLogs,

		HeaderActionAWS,
		ControlActionTerminate,
	})

	// NO-ACTIONS / HEADERS
	switch action {
	case HeaderActionAWS,
		HeaderActionJamulus,
		HeaderActionSCP:
		m.DisplayControlInstance(instance)
		return
	}

	// AWS-ACTIONS
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

	// RUNNING-ACTIONS
	if instance.State == nil || instance.State.Name == nil || *instance.State.Name != "running" {
		fmt.Println(common.AWSPrefix(), "This action requires the instance to be running. Starting ...")

		// start instance
		var err error
		if _, err = m.ec.StartInstances(&ec2.StartInstancesInput{
			InstanceIds: []*string{instance.InstanceId},
		}); err != nil {
			panic(err)
		}

		// wait for instance
		if instance, err = waiter.WaitForState(m.ec, *instance.InstanceId, "running"); err != nil {
			panic(err)
		}
	}

	// SSH-ACTIONS
	fmt.Println(common.SSHPrefix(), "This action requires SSH")
	ssh, err := waiter.WaitForSSHInstance(instance, tpl)
	if err != nil {
		panic(err)
	}
	switch action {
	case ControlActionStartJamulus:
		ctl.StartJamulus(ssh, tpl)
		return

	case ControlActionStopJamulus:
		ctl.StopJamulus(ssh, true)
		return

	case ControlActionToggleRecording:
		ctl.JamulusRecord(ssh, ctl.JamulusRecModeToggle, true)
	}
}
