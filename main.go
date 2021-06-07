package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/menu"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"log"
	"os"
	"path"

	// apply windows patch
	_ "github.com/darmiel/jamulus-aws-deploy/internal/thin"
)

const (
	Region = "eu-central-1"
)

func main() {
	fmt.Println("aws-deploy - compiled for", tsess.Owner)

	// create data/templates dir
	_ = os.MkdirAll(path.Join("data", "templates"), os.ModePerm)
	_ = os.MkdirAll(path.Join("data", "keys"), os.ModePerm)

	// create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		log.Fatalln("Error creating session:", err)
		return
	}
	ec := ec2.New(sess, aws.NewConfig().WithRegion(Region))

	m := menu.NewMenu(ec)
	m.DisplayListInstances()

	// catch errors
	if a := recover(); a != nil {
		fmt.Println(common.ERRPrefix(), "JAWS crashed:", a)
		fmt.Println(common.ERRPrefix(), "Please try again.")
	}
}
