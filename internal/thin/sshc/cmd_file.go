package sshc

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	lsRegex = regexp.MustCompile(
		`(?m)^([drwx-]{10})\s+([1-9])\s+([A-Za-z0-9-_]+)\s+([A-Za-z0-9-_]+)\s+([0-9]+)\s+([0-9-]+)\s+([0-9:.]+)\s+([-+0-9]+)\s+(.*)$`)
)

func (s *SSHC) DownloadFileBase64(lpath, rpath string) (int, error) {
	resp := s.MustExecute("sudo cat %s | base64", strconv.Quote(rpath))
	if resp.StatusCode != 0 {
		return -1, errors.New("file probably not found, or missing permissions")
	}
	out := make([]byte, len(resp.Data))
	n, err := base64.StdEncoding.Decode(out, resp.Data)
	if err != nil {
		return -1, err
	}
	return n, os.WriteFile(lpath, out[0:n], 0644)
}

type FileInfo struct {
	Permission string
	NumLinks   uint8
	Owner      string
	Group      string
	Size       uint64
	Time       time.Time
	FileName   string
}

func (i *FileInfo) IsDir() bool {
	return strings.HasPrefix(i.Permission, "d")
}

func (s *SSHC) DirLS(dir string) []*FileInfo {
	resp := s.MustExecute("sudo ls -la --time-style=full %s", strconv.Quote(dir))
	res := make([]*FileInfo, 0)
	if resp.StatusCode != 0 {
		return res
	}
	for _, m := range lsRegex.FindAllStringSubmatch(string(resp.Data), -1) {
		name := m[9]
		if name == "." || name == ".." {
			continue
		}

		numLinks, err := strconv.ParseInt(m[2], 10, 8)
		if err != nil {
			fmt.Println("WARN (numLinks) ::", err)
			continue
		}
		size, err := strconv.ParseInt(m[5], 10, 64)
		if err != nil {
			fmt.Println("WARN (size) ::", err)
			continue
		}

		t, err := time.Parse("2006-01-02 15:04:05.000000000 -0700",
			fmt.Sprintf("%s %s %s", m[6], m[7], m[8]))
		if err != nil {
			fmt.Println("WRN:", "error parsing date:", err)
			continue
		}

		i := &FileInfo{
			Permission: m[1],
			NumLinks:   uint8(numLinks),
			Owner:      m[3],
			Group:      m[4],
			Size:       uint64(size),
			FileName:   name,
			Time:       t,
		}
		res = append(res, i)
	}
	return res
}

///

func (s *SSHC) DirExists(dir string) bool {
	return s.IsStatusCode(0, "[ -d %s ]", strconv.Quote(dir))
}

func (s *SSHC) DirCreate(dir string) bool {
	if s.DirExists(dir) {
		return false
	}
	s.MustExecute("mkdir -p %s", strconv.Quote(dir))
	return true
}

///

func (s *SSHC) FileExists(file string) bool {
	return s.IsStatusCode(0, "[ -f %s ]", strconv.Quote(file))
}
