package templates

import (
	"time"
)

type Template struct {
	LocalTemplate string            `json:"-"`
	Template      *TemplateInfo     `json:"Template"`
	Instance      *TemplateInstance `json:"Instance"`
	Jamulus       *TemplateJamulus  `json:"Jamulus"`
}

///////////

type TemplateInfo struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Author      *TemplateAuthor `json:"Author"`
	Date        time.Time       `json:"Date"`
}

type TemplateAuthor struct {
	Name    string `json:"Name"`
	Contact string `json:"Contact"`
}

//

type TemplateInstance struct {
	Type              string           `json:"Type"`
	SecurityGroupName string           `json:"Security Group Name"` // ok
	KeyPair           *TemplateKeyPair `json:"Key Pair"`
	AMI               string           `json:"AMI"`
}

type TemplateKeyPair struct {
	Name      string `json:"Name"`
	LocalPath string `json:"Local Path"`
}

//

type TemplateJamulus struct {
	MaxUsers             uint                      `json:"Max Users"`
	WelcomeMessage       string                    `json:"Welcome Message"`
	Public               *TemplateJamulusPublic    `json:"Public"`
	Recording            *TemplateJamulusRecording `json:"Recording"`
	FastUpdate           bool                      `json:"Fast Update"`
	EnableMultiThreading bool                      `json:"Enable Multithreading"`
	LogPath              string                    `json:"Log Path"`
}

type TemplateJamulusPublic struct {
	CentralServer string                           `json:"Central Server"`
	ServerInfo    *TemplateJamulusPublicServerInfo `json:"Server Info"`
}

type TemplateJamulusPublicServerInfo struct {
	Name    string `json:"Name"`
	City    string `json:"City"`
	Country string `json:"Country"`
}

type TemplateJamulusRecording struct {
	Path       string `json:"Path"`
	AutoRecord bool   `json:"Auto Record"`
}
