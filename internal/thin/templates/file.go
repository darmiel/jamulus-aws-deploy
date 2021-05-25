package templates

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const (
	PermMode = 0755
)

var (
	TemplateDir = path.Join("data", "templates")
)

// read / write
func FromFile(name string) (tpl *Template, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path.Join(TemplateDir, name)); err != nil {
		return
	}
	if err = json.Unmarshal(data, &tpl); err == nil {
		tpl.LocalTemplate = name
	}
	return
}

func Must(tpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func (t *Template) ToFile(name string) (err error) {
	var data []byte
	if data, err = json.Marshal(t); err != nil {
		return
	}
	err = os.WriteFile(path.Join(TemplateDir, name), data, PermMode)
	return
}
