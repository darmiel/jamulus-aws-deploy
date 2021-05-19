package tpl

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type CreateInstanceTemplate struct {
	TemplateName        string
	TemplateDescription string
	//
	InstanceType    string
	KeyPair         string
	KeyPairPath     string
	SecurityGroupID string
	//
	IsTemplate bool `json:"-"`
}

const (
	NoTemplate = "🤷 No Template / New Template"
)

func SelectTemplate() *CreateInstanceTemplate {
	dir := "templates"
	// check if template folder exists
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		// create dir
		// TODO: change me (0755)
		if err := os.Mkdir("templates", 0755); err != nil {
			log.Fatalln("Error creating templates folder:", err)
			return nil
		}
	}

	glob, err := filepath.Glob("templates/*.json")
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
		res.IsTemplate = true

		i++

		name := fmt.Sprintf("%d. %s (%s)", i, res.TemplateName, res.TemplateDescription)
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
