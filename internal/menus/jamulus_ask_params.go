package menus

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/darmiel/jamulus-aws-deploy/internal/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/tpl"
	"log"
	"strconv"
)

type AskJamulusParamsMenu *Menu

func NewAskJamulusParamsMenu(parent *Menu) AskJamulusParamsMenu {
	menu := &Menu{
		Parent: parent,
	}
	menu.Print = func() {
		temp := tpl.SelectTemplate(tpl.TemplateTypeJamulus)
		if temp == nil {
			log.Fatalln("error selecting template")
			return
		}

		// make public?
		if temp.Jamulus.CentralServer == nil {
			if common.Bool("Make Server Public", false) {
				// central server
				temp.Jamulus.CentralServer = aws.String(common.FlatSelect("Select Server", map[string]string{
					"Any Genre 1":             "anygenre1.jamulus.io:22124",
					"Any Genre 2":             "anygenre2.jamulus.io:22224",
					"Any Genre 3":             "anygenre3.jamulus.io:22624",
					"Genre Rock":              "rock.jamulus.io:22424",
					"Genre Jazz":              "jazz.jamulus.io:22324",
					"Genre Classical/Folk":    "classical.jamulus.io:22524",
					"Genre Choral/Barbershop": "choral.jamulus.io:22724",
				}))

				// server info
				data := struct {
					Name   string
					City   string
					Locale string
				}{}
				q := []*survey.Question{
					{
						Name:     "Name",
						Prompt:   &survey.Input{Message: "Public Server-Name"},
						Validate: survey.Required,
					},
					{
						Name:     "City",
						Prompt:   &survey.Input{Message: "Public Server-City"},
						Validate: survey.Required,
					},
					{
						Name:     "Locale",
						Prompt:   &survey.Input{Message: "Public Server-Locale"},
						Validate: survey.Required,
					},
				}

				if err := survey.Ask(q, &data); err != nil {
					log.Fatalln("error reading your answer:", err)
					return
				}

				temp.Jamulus.ServerInfo = aws.String(fmt.Sprintf("%s;%s;%s", data.Name, data.City, data.Locale))
			} else {
				temp.Jamulus.CentralServer = aws.String("-")
			}
		}

		// max users
		if temp.Jamulus.MaxUsers == nil {
			usersStr := common.StringValidate("Max Users", "25", func(ans interface{}) error {
				s, ok := ans.(string)
				if !ok {
					return errors.New("not a string")
				}
				_, err := strconv.Atoi(s)
				return err
			})
			atoi, err := strconv.Atoi(usersStr)
			if err != nil {
				log.Fatalln("error reading max users:", err)
				return
			}
			temp.Jamulus.MaxUsers = &atoi
		}

		// welcome message
		if temp.Jamulus.WelcomeMessage == nil {
			if temp.Jamulus.WelcomeMessage = aws.String(common.String("Welcome Message", "Hallo!")); *temp.Jamulus.WelcomeMessage == "" {
				temp.Jamulus.WelcomeMessage = aws.String("-")
			}
		}

		// fast update
		if temp.Jamulus.FastUpdate == nil {
			temp.Jamulus.FastUpdate = aws.Bool(common.Bool("Use Fast Update?", true))
		}

		// multi threading
		if temp.Jamulus.EnableMultiThreading == nil {
			temp.Jamulus.EnableMultiThreading = aws.Bool(common.Bool("Enable Multi Threading", true))
		}

		// recording
		if temp.Jamulus.RecordingPath == nil {
			if common.Bool("Enable Recording", true) {
				temp.Jamulus.RecordingPath = aws.String(common.String("Recording Path", "/tmp/recordings"))
				temp.Jamulus.NoRecordOnStart = aws.Bool(common.Bool("No Record On Start", true))
			} else {
				temp.Jamulus.RecordingPath = aws.String("-")
				temp.Jamulus.NoRecordOnStart = aws.Bool(true)
			}
		}

		// logging
		if temp.Jamulus.LogPath == nil {
			if common.Bool("Enable Logging", true) {
				temp.Jamulus.LogPath = aws.String(common.String("Log Path", "/tmp/logs"))
			} else {
				temp.Jamulus.LogPath = aws.String("-")
			}
		}

		temp.AskSave()

		log.Printf("Settings: %+v", temp.Jamulus)
	}
	return menu
}
