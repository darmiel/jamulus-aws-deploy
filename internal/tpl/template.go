package tpl

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"path"
)

const (
	TemplateTypeInstance = iota
	TemplateTypeJamulus
)

type CreateInstanceTemplate struct {
	Template struct {
		TemplateName        string
		TemplateDescription string
		TemplateType        uint64
		IsTemplate          bool `json:"-"`
	}
	Instance struct {
		InstanceType    string
		KeyPair         string
		KeyPairPath     string
		SecurityGroupID string
	}
	Jamulus struct {
		CentralServer *string // ok
		ServerInfo    *string // ok
		FastUpdate    *bool   // ok
		LogPath       *string // ok
		// HTMLStatusFile       *string //
		RecordingPath        *string // ok
		NoRecordOnStart      *bool   // ok
		EnableMultiThreading *bool   // ok
		MaxUsers             *int    // ok
		WelcomeMessage       *string // ok
	}
}

func (t *CreateInstanceTemplate) AskSave() {
	// save template
	if !t.Template.IsTemplate {
		var saveTemplate bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Save [Server] Template?",
			Default: false,
		}, &saveTemplate); err != nil {
			log.Fatalln("Error reading your answer:", err)
			return
		}

		// ask for name and description
		if err := survey.Ask([]*survey.Question{
			{
				Name:     "TemplateName",
				Prompt:   &survey.Input{Message: "Template Name"},
				Validate: survey.Required,
			},
			{
				Name:   "TemplateDescription",
				Prompt: &survey.Input{Message: "Template Description"},
			},
		}, &t.Template); err != nil {
			log.Fatalln("Error reading your answer:", err)
			return
		}

		// encode to json
		data, err := json.Marshal(t)
		if err != nil {
			log.Fatalln("Error encoding to JSON:", err)
			return
		}

		// generate uuid
		b := make([]byte, 16)
		if _, err = rand.Read(b); err != nil {
			log.Fatalln("Error generating id:", err)
			return
		}
		uuid := fmt.Sprintf("%x-%x-%x-%x-%x.json", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		// write
		if err := os.WriteFile(path.Join("templates", uuid), data, 0755); err != nil {
			log.Fatalln("Error writing to file:", err)
			return
		}
	}
}
