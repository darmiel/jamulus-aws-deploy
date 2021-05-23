package sshc

import "strconv"

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
