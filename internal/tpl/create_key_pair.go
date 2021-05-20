package tpl

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
	"path/filepath"
)

func (t *CreateInstanceTemplate) CreateKeyPair(ec *ec2.EC2) {
	pair, err := ec.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(t.Instance.KeyPair),
	})
	if err != nil {
		log.Fatalln("Error creating key pair:", err)
		return
	}

	// save key
	var outPath string
	if t.Instance.KeyPairPath != "" {
		if _, err := os.Stat(t.Instance.KeyPairPath); os.IsNotExist(err) {
			outPath = t.Instance.KeyPairPath
		}
	}

	// ask for path
	if outPath == "" {
		if err := survey.AskOne(&survey.Input{
			Message: "Path of your key pair",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}, &outPath); err != nil {
			log.Fatalln("Error reading your answer:", err)
			return
		}
	}

	// save to path
	if err := os.WriteFile(outPath, []byte(*pair.KeyMaterial), 0644); err != nil {
		log.Fatalln("Error writing cert:", err)
		return
	}
}
