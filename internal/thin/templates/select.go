package templates

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	Permission = 0777
)

func SelectTemplate() *Template {
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

	opts := make(map[string]*Template)
	optsStr := make([]string, 0)

	i := 0
	for _, file := range glob {
		if strings.Contains(file, "/") {
			file = file[strings.LastIndex(file, "/")+1:]
		}

		res, err := FromFile(file)
		if err != nil {
			panic(err)
		}

		i++
		name := fmt.Sprintf("%d. %s (%s)", i, res.Template.Name, res.Template.Description)
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

	return opts[resp]
}
