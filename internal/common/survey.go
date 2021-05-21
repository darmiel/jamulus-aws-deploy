package common

import (
	"github.com/AlecAivazis/survey/v2"
	"log"
)

func Bool(message string, def bool) (resp bool) {
	q := &survey.Confirm{
		Message: message,
		Default: def,
	}
	if err := survey.AskOne(q, &resp); err != nil {
		log.Fatalln("error reading your answer:", err)
	}
	return
}

func String(message, def string) (resp string) {
	q := &survey.Input{
		Message: message,
		Default: def,
	}
	if err := survey.AskOne(q, &resp); err != nil {
		log.Fatalln("error reading your answer:", err)
	}
	return
}

func StringValidate(message, def string, val survey.Validator) string {
	data := struct {
		Resp string
	}{}
	q := []*survey.Question{
		{
			Name: "Resp",
			Prompt: &survey.Input{
				Message: message,
				Default: def,
			},
			Validate: val,
		},
	}
	if err := survey.Ask(q, &data); err != nil {
		log.Fatalln("error reading your answer:", err)
	}
	return data.Resp
}

func Select(message string, opts []string) (resp string) {
	q := &survey.Select{
		Message: message,
		Options: opts,
	}
	if err := survey.AskOne(q, &resp); err != nil {
		log.Fatalln("error reading your answer:", err)
	}
	return
}

func FlatSelect(message string, optMap map[string]string) (resp string) {
	opts := make([]string, len(optMap))
	i := 0
	for k, _ := range optMap {
		opts[i] = k
		i++
	}
	return optMap[Select("Select Server", opts)]
}