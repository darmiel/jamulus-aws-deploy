package templates

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const PermMode = 0755

// read / write
func FromFile(path string) (tpl *Template, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	err = json.Unmarshal(data, &tpl)
	return
}

func Must(tpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

func (t *Template) ToFile(path string) (err error) {
	var data []byte
	if data, err = json.Marshal(t); err != nil {
		return
	}
	err = os.WriteFile(path, data, PermMode)
	return
}
