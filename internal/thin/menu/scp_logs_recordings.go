package menu

import (
	"fmt"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/sshc"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/templates"
	"github.com/dustin/go-humanize"
	"path/filepath"
	"strings"
)

func dlPath(rpath string, ssh *sshc.SSHC, tpl *templates.Template) string {
	kp := tpl.Instance.KeyPair.LocalPath

	var builder strings.Builder
	builder.WriteString(common.Color("scp", "#FFFFFF").String())
	builder.WriteRune(' ')

	abs, err := filepath.Abs(kp)
	if err != nil {
		return ""
	}

	builder.WriteString(common.Color(fmt.Sprintf("-r -i %s", abs), "#D290E4").String())
	builder.WriteRune(' ')
	builder.WriteString(common.Color("ec2-user", "#e67e22").String())
	builder.WriteRune('@')
	builder.WriteString(common.Color(ssh.Client().Config.Addr, "#2ecc71").String())
	builder.WriteRune(':')
	builder.WriteString(common.Color(rpath, "#1abc9c").String())
	builder.WriteRune(' ')

	if abs, err = filepath.Abs("data"); err != nil {
		return ""
	}

	builder.WriteString(common.Color(abs, "#2980b9").String())
	return builder.String()
}

func (m *Menu) ListLogs(ssh *sshc.SSHC, tpl *templates.Template) {
	if tpl.Jamulus.LogPath == "" {
		fmt.Println(common.ERRPrefix(), "Logging is not enabled.")
		return
	}

	ls := ssh.DirLS(tpl.Jamulus.LogPath)

	fmt.Println()
	fmt.Println(common.SSHPrefix(), "Found",
		common.Color(humanize.Comma(int64(len(ls))), "#71BEF2"), "files in log-dir:")
	fmt.Println(common.SSHPrefix(), "Download-CMD:")
	fmt.Println(common.SSHPrefix(), dlPath(tpl.Jamulus.LogPath, ssh, tpl))

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
	fmt.Println(common.SSHPrefix(), "Download-CMD:")
	fmt.Println(common.SSHPrefix(), dlPath(tpl.Jamulus.Recording.Path, ssh, tpl))

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
