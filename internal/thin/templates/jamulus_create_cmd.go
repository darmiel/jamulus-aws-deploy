package templates

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	JamulusDockerImage = "grundic/jamulus"
)

func (t *TemplateJamulus) CreateArgs() string {
	args := make([]string, 0)

	// sudo
	args = append(args, "sudo", "docker", "run", "-d", "--rm")

	// ports
	args = append(args, "-p",
		fmt.Sprintf("%d:%d/udp", 22124, 22124))

	// volumes
	if c := t.LogPath; c != "" {
		q := strconv.Quote(c)
		args = append(args, fmt.Sprintf("-v %s:%s", q, q))
	}
	if c := t.Recording.Path; c != "" {
		q := strconv.Quote(c)
		args = append(args, fmt.Sprintf("-v %s:%s", q, q))
	}

	// docker image
	args = append(args, JamulusDockerImage)

	// default params
	args = append(args, "-n -s -p 22124")

	// append params
	// central server
	if t.Public.CentralServer != "" {
		args = append(args, t.Public.CreateArgs())
	}

	args = append(args, "--numchannels",
		strconv.FormatInt(int64(t.MaxUsers), 10))

	// fast update
	if t.FastUpdate {
		args = append(args, "--fastupdate")
	}

	// log path
	if c := t.LogPath; c != "" {
		args = append(args, "--log", strconv.Quote(c))
	}

	// recording
	if c := t.Recording.Path; c != "" {
		args = append(args, "--recording", strconv.Quote(c))
		if !t.Recording.AutoRecord {
			args = append(args, "--norecord")
		}
	}

	// multithreading
	if t.EnableMultiThreading {
		args = append(args, "--multithreading")
	}

	// welcome message
	if c := t.WelcomeMessage; c != "" {
		args = append(args, "--welcomemessage", strconv.Quote(c))
	}

	return strings.Join(args, " ")
}

func (p *TemplateJamulusPublic) CreateArgs() string {
	args := make([]string, 0)
	args = append(args, "--centralserver",
		p.CentralServer,
		"--serverinfo",
		strconv.Quote(p.ServerInfo.String()))
	return strings.Join(args, " ")
}

func (i *TemplateJamulusPublicServerInfo) String() string {
	return fmt.Sprintf("%s;%s;%s", i.Name, i.City, i.Country)
}
