package tpl

import (
	"github.com/briandowns/spinner"
	"time"
)

func NewSpinner(prefix, final string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Prefix = prefix + " "
	s.FinalMSG = final
	s.Start()
	return s
}
