package common

import (
	"github.com/aws/aws-sdk-go/service/ec2"
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
