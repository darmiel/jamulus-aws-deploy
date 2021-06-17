package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/awswrap"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/menu"
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
	fmt.Println("⭐️ aws-deploy - compiled for", common.Owner)

	/// FS

	// create data/templates dir
	_ = os.MkdirAll(path.Join("data", "templates"), os.ModePerm)
	_ = os.MkdirAll(path.Join("data", "keys"), os.ModePerm)

	// check if credential file exists
	var awsCredPath = defaults.SharedCredentialsFilename()
	if file, err := os.Stat(awsCredPath); os.IsNotExist(err) || file == nil || file.IsDir() {
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
		st     *sts.STS
	)

	for {
		// create session
		if sess, err = session.NewSession(awsCfg); err != nil {
			log.Fatalln("Error creating session:", err)
			return
		}

		// STS - AWS Security Token Service
		// -> Check Credentials
		st = sts.New(sess)
		if caller, err = st.GetCallerIdentity(&sts.GetCallerIdentityInput{}); err != nil {
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
		common.Color(common.DeRef(caller.UserId), "66C2CD"))

	// EC2 - Elastic Compute Cloud
	ec := ec2.New(sess)

	wrap := &awswrap.AWSWrap{
		CredPath: awsCredPath,
		Region:   Region,
		Config:   awsCfg,
		STS:      st,
		EC2:      ec,
	}

	m := menu.NewMenu(ec, wrap)
	m.DisplayListInstances(common.Owner, false, true)

	// catch errors
	if a := recover(); a != nil {
		fmt.Println(common.ERRPrefix(), "JAWS crashed:", a)
		fmt.Println(common.ERRPrefix(), "Please try again.")
	}
}
