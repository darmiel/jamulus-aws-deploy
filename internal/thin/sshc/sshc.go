package sshc

import (
	"github.com/melbahja/goph"
	"log"
)

type SSHC struct {
	client *goph.Client
}

func Must(g *goph.Client, err error) *SSHC {
	if err != nil {
		log.Fatalln("error connecting to ssh:", err)
		return nil
	}
	return &SSHC{g}
}
