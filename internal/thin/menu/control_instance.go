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
	ControlActionGoBack          = "👋 | Go Back"
	ControlActionStartJamulus    = "JAM | 🚀 | Start Jamulus"
	ControlActionStopJamulus     = "JAM | 🔻 | Stop Jamulus"
	ControlActionToggleRecording = "JAM | 🎤 | Toggle Recording"
	ControlActionGetRecordings   = "SCP | 📂 | Browse Recordings"
	ControlActionGetLogs         = "SCP | 📂 | Browse Logs"
	ControlActionTerminate       = "AWS | 🚮 | Terminate"
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
		ControlActionGoBack,
		ControlActionStartJamulus,
		ControlActionStopJamulus,
		ControlActionToggleRecording,
		ControlActionGetRecordings,
		ControlActionGetLogs,
		ControlActionTerminate,
	})

	// OTHER ACTIONS
	if action == ControlActionGoBack {
		return
	}

	// AWS-ACTIONS
	if action == ControlActionTerminate {
		// confirm
		if common.Bool("Terminate #"+*instance.InstanceId+"?", false) {
			_, err := m.ec.TerminateInstances(&ec2.TerminateInstancesInput{
				InstanceIds: []*string{instance.InstanceId},
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(common.AWSPrefix(), "Terminating", *instance.InstanceId, "...")
		}
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

	case ControlActionStopJamulus:
		ctl.StopJamulus(ssh, true)

	case ControlActionToggleRecording:
		ctl.JamulusRecord(ssh, ctl.JamulusRecModeToggle, true)

	case ControlActionGetLogs:
		m.ListLogs(ssh, tpl)

	case ControlActionGetRecordings:
		m.ListRecordings(ssh, tpl)
	}

	// go back to control instance
	m.DisplayControlInstance(instance)
}
