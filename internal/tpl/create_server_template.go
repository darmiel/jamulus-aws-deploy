package tpl

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	DefaultKeyPair = "jamulus-cert"
	NoTemplate     = "🤷 No Template / New Template"
	TemplateDir    = "templates"
	Permission     = 0755
)

func SelectTemplate(ty uint64) *CreateInstanceTemplate {
	// check if template folder exists
	if info, err := os.Stat(TemplateDir); err != nil || !info.IsDir() {
		// create dir
		if err := os.Mkdir(TemplateDir, Permission); err != nil {
			log.Fatalln("Error creating templates folder:", err)
			return nil
		}
	}

	glob, err := filepath.Glob(path.Join(TemplateDir, "*.json"))
	if err != nil {
		log.Fatalln("Error listing files:", err)
		return nil
	}

	opts := make(map[string]*CreateInstanceTemplate)
	optsStr := make([]string, 1)
	optsStr[0] = NoTemplate

	i := 0
	for _, file := range glob {
		// read data
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("WARN :: Error reading template", file, ":", err)
			continue
		}

		// parse file
		res := new(CreateInstanceTemplate)
		if err := json.Unmarshal(data, res); err != nil {
			log.Println("WARN :: Error parsing template", file, ":", err)
			continue
		}
		res.Template.IsTemplate = true

		// check template type
		if res.Template.TemplateType != ty {
			continue
		}

		i++
		name := fmt.Sprintf("%d. %s (%s)", i, res.Template.TemplateName, res.Template.TemplateDescription)
		opts[name] = res
		optsStr = append(optsStr, name)
	}

	q := &survey.Select{
		Message: "Select Template",
		Options: optsStr,
	}
	var resp string
	if err := survey.AskOne(q, &resp); err != nil {
		log.Fatalln("Error reading answer:", err)
		return nil
	}

	if resp == NoTemplate {
		return &CreateInstanceTemplate{}
	}

	return opts[resp]
}
