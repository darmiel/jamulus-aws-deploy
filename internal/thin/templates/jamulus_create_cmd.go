package templates

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *TemplateJamulus) CreateArgs(sudo, docker bool, name string) string {
	var builder strings.Builder

	// append docker
	if docker {
		if sudo {
			builder.WriteString("sudo ")
		}
		builder.WriteString(`docker run -d --rm --name "`)
		builder.WriteString(name)
		builder.WriteString(`" `)
		builder.WriteString("grundic/jamulus ")
	}

	// append params
	// central server
	if t.Public.CentralServer != "" {
		builder.WriteString(t.Public.CreateArgs())
		builder.WriteRune(' ')
	}

	// fast update
	if t.FastUpdate {
		builder.WriteString("--fastupdate ")
	}

	// log path
	if t.LogPath != "" {
		builder.WriteString(`--log "`)
		builder.WriteString(t.LogPath)
		builder.WriteString(`" `)
	}

	// recording
	if t.Recording.Path != "" {
		builder.WriteString(`--recording "`)
		builder.WriteString(t.Recording.Path)
		builder.WriteString(`" `)
		if !t.Recording.AutoRecord {
			builder.WriteString("--norecord ")
		}
	}

	// multithreading
	if t.EnableMultiThreading {
		builder.WriteString("--multithreading ")
	}

	// num channels
	builder.WriteString("--numchannels ")
	builder.WriteString(strconv.FormatInt(int64(t.MaxUsers), 10))
	builder.WriteRune(' ')

	// welcome message
	if t.WelcomeMessage != "" {
		builder.WriteString(`--welcomemessage "`)
		builder.WriteString(t.WelcomeMessage)
		builder.WriteString(`" `)
	}

	return builder.String()
}

func (p *TemplateJamulusPublic) CreateArgs() string {
	var builder strings.Builder

	// central server
	builder.WriteString("--centralserver ")
	builder.WriteString(p.CentralServer)
	builder.WriteRune(' ')

	// server info
	builder.WriteString("--serverinfo ")
	builder.WriteString(p.ServerInfo.String())

	return builder.String()
}

func (i *TemplateJamulusPublicServerInfo) String() string {
	return fmt.Sprintf("%s;%s;%s", i.Name, i.City, i.Country)
}
