package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/menu"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/tsess"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/wizard"
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

	/// FS

	// create data/templates dir
	_ = os.MkdirAll(path.Join("data", "templates"), os.ModePerm)
	_ = os.MkdirAll(path.Join("data", "keys"), os.ModePerm)

	// check if credential file exists
	var awsCredPath = defaults.SharedCredentialsFilename()
	if file, err := os.Stat(awsCredPath); os.IsNotExist(err) || file.IsDir() {
		// Start AWS wizard
		if err := wizard.StartAWSCredWizard(); err != nil {
			fmt.Println(common.ERRPrefix(), err)
			return
		}
	}

	//// AWS
	// AWS Config
	awsCfg := aws.NewConfig().WithRegion(Region)

	var (
		err    error
		sess   *session.Session
		caller *sts.GetCallerIdentityOutput
	)

	for {
		// create session
		if sess, err = session.NewSession(awsCfg); err != nil {
			log.Fatalln("Error creating session:", err)
			return
		}

		// STS - AWS Security Token Service
		// -> Check Credentials
		st := sts.New(sess)
		if caller, err = st.GetCallerIdentity(&sts.GetCallerIdentityInput{}); err != nil {
			fmt.Println()
			fmt.Println(common.ERRPrefix(), "Unknown identity")
			fmt.Println(common.ERRPrefix(), "It is most likely that your login details are not correct.")
			if !common.Bool("Do you want to enter new login data?", true) {
				return
			}
			if err = wizard.StartAWSCredWizard(); err != nil {
				fmt.Println(common.ERRPrefix(), "Error reading your new login data:", err)
			}
			continue
		}

		break
	}

	// say hello 👋
	fmt.Println(common.AWSPrefix(), "Logged in with user",
		common.Color(common.Deref(caller.UserId), "66C2CD"))

	// EC2 - Elastic Compute Cloud
	ec := ec2.New(sess)

	m := menu.NewMenu(ec)
	m.DisplayListInstances()

	// catch errors
	if a := recover(); a != nil {
		fmt.Println(common.ERRPrefix(), "JAWS crashed:", a)
		fmt.Println(common.ERRPrefix(), "Please try again.")
	}
}
