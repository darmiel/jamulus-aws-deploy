package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/muesli/termenv"
	"strings"
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

var p = termenv.ColorProfile()

func Profile() termenv.Profile {
	return p
}

func Color(msg, color string) termenv.Style {
	if !strings.HasPrefix(color, "#") {
		color = "#" + color
	}
	return termenv.String(msg).Foreground(Profile().Color(color))
}

func SSHPrefix() termenv.Style {
	return termenv.String(" SSH ").Foreground(p.Color("0")).Background(p.Color("#3498db"))
}

func ERRPrefix() termenv.Style {
	return termenv.String(" ERR ").Foreground(p.Color("0")).Background(p.Color("#E88388"))
}

func JAMPrefix() termenv.Style {
	return termenv.String(" JAM ").Foreground(p.Color("0")).Background(p.Color("#D290E4"))
}

func AWSPrefix() termenv.Style {
	return termenv.String(" AWS ").Foreground(p.Color("0")).Background(p.Color("#DBAB79"))
}

func LvlPrint(prefix interface{}, message string) {
	for _, line := range strings.Split(message, "\n") {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}
		fmt.Println(prefix, line)
	}
}

type PrintReport struct{}

func (*PrintReport) Report(resp interface{}) {
	var st string
	switch t := resp.(type) {
	case *sshc.SSHCCommandResult:
		if t.StatusCode == 0 {
			st = "👍"
		} else {
			st = "🤬"
		}
		fmt.Println(st)
		if t.StatusCode != 0 {
			for _, line := range strings.Split(string(t.Data), "\n") {
				if len(strings.TrimSpace(line)) <= 0 {
					continue
				}
				fmt.Println(ERRPrefix(), termenv.String(line).Foreground(p.Color("#DBAB79")))
			}
		}
	case bool:
		if t {
			st = "👍"
		} else {
			st = "🤬"
		}
		fmt.Println(st)
	default:
		fmt.Printf("Report :: invalid type: %T (%v)\n", t, t)
	}
}

func PrintState(prefix interface{}, msg string) *PrintReport {
	fmt.Print(prefix, " 🔨 | ", msg, " ... ")
	return &PrintReport{}
}

func PrintSSHState(msg string) *PrintReport {
	return PrintState(SSHPrefix(), msg)
}

func DeRef(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func EmptyDef(str string, def string) string {
	if str != "" {
		return str
	}
	return def
}
