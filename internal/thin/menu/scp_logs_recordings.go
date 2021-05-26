package menu

import (
	"fmt"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/dustin/go-humanize"
)

func (m *Menu) ListLogs(ssh *sshc.SSHC, tpl *templates.Template) {
	if tpl.Jamulus.LogPath == "" {
		fmt.Println(common.ERRPrefix(), "Logging is not enabled.")
		return
	}

	ls := ssh.DirLS(tpl.Jamulus.LogPath)

	fmt.Println()
	fmt.Println(common.SSHPrefix(), "Found",
		common.Color(humanize.Comma(int64(len(ls))), "#71BEF2"), "files in log-dir:")

	for _, info := range ls {
		var icon string
		if info.IsDir() {
			icon = "📂"
		} else {
			icon = "📝"
		}

		fmt.Println(common.SSHPrefix(),
			icon,
			">",
			info.FileName,
			common.Color(humanize.Bytes(info.Size), "#DBAB79"),
			common.Color(info.Time.Format("02.01.2006 15:04:05"), "#D290E4"))
	}
	fmt.Println()
}

func (m *Menu) ListRecordings(ssh *sshc.SSHC, tpl *templates.Template) {
	if tpl.Jamulus.Recording == nil || tpl.Jamulus.Recording.Path == "" {
		fmt.Println(common.ERRPrefix(), "Recording is not enabled.")
		return
	}

	ls := ssh.DirLS(tpl.Jamulus.Recording.Path)

	fmt.Println()
	fmt.Println(common.SSHPrefix(), "Found",
		common.Color(humanize.Comma(int64(len(ls))), "#71BEF2"), "files in rec-dir:")

	for _, info := range ls {
		var icon string
		if info.IsDir() {
			icon = "📂"
		} else {
			icon = "📝"
		}

		fmt.Println(common.SSHPrefix(),
			icon,
			">",
			info.FileName,
			common.Color(humanize.Bytes(info.Size), "#DBAB79"),
			common.Color(info.Time.Format("02.01.2006 15:04:05"), "#D290E4"))
	}
	fmt.Println()
}
