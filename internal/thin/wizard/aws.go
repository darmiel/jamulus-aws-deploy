package wizard

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/darmiel/jamulus-aws-deploy/internal/thin/common"
	"os"
	"strings"
)

var qs = []*survey.Question{
	{
		Name:     "id",
		Prompt:   &survey.Input{Message: "Enter Key ID"},
		Validate: survey.Required,
	},
	{
		Name:     "key",
		Prompt:   &survey.Input{Message: "Enter Access Key"},
		Validate: survey.Required,
	},
}

func askAWSCred() (err error) {
	data := struct {
		KeyID     string `survey:"id"`
		AccessKey string `survey:"key"`
	}{}
	if err = survey.Ask(qs, &data); err != nil {
		return
	}

	// clean input
	data.KeyID = strings.TrimSpace(data.KeyID)
	data.AccessKey = strings.TrimSpace(data.AccessKey)

	// generate file
	content := fmt.Sprintf(`[default]%saws_access_key_id = "%s"%saws_secret_access_key = "%s"`,
		"\n", data.KeyID, "\n", data.AccessKey)

	credPath := defaults.SharedCredentialsFilename()

	awsCredDir, awsCredFile := common.SplitPath(credPath)
	fmt.Println(common.AWSPrefix(), "Cred-Dir:", awsCredDir)
	fmt.Println(common.AWSPrefix(), "Cred-File:", awsCredFile)

	// mkdirs
	_ = os.MkdirAll(awsCredDir, os.ModePerm)

	// write to file
	err = os.WriteFile(credPath, []byte(content), os.ModePerm)
	return
}

func showAWSAccessKeyObtain(sc *bufio.Scanner) {
	fmt.Println(common.AWSPrefix(), "Visit",
		common.Color("https://console.aws.amazon.com/iam/home#/security_credentials", "66C2CD"))
	fmt.Println(common.AWSPrefix(), "< press", common.Color("enter", "#A8CC8C"), "to continue >")

	sc.Scan()
	fmt.Println(common.AWSPrefix(), "Open menu",
		common.Color("Access Keys (access key ID and secret access key)", "66C2CD"),
		"(currently the 3rd menu)")
	fmt.Println(common.AWSPrefix(), "< press", common.Color("enter", "#A8CC8C"), "to continue >")

	sc.Scan()
	fmt.Println(common.AWSPrefix(), "Click on",
		common.Color("Create New Access Key", "66C2CD"))
	fmt.Println(common.AWSPrefix(), "< press", common.Color("enter", "#A8CC8C"), "to continue >")

	sc.Scan()
	fmt.Println(common.AWSPrefix(), "You should now see",
		common.Color("Access Key ID", "A8CC8C"),
		"and",
		common.Color("Secret Access Key", "A8CC8C"))
	fmt.Println(common.AWSPrefix(), "< press", common.Color("enter", "#A8CC8C"), "to continue >")
}

func StartAWSCredWizard() (err error) {
	fmt.Println(common.AWSPrefix(), "It looks like your AWS account is not configured yet.")
	fmt.Println(common.AWSPrefix(), "The AWS account is needed to manage instances.")
	fmt.Println(common.AWSPrefix(), "Do you want to store your login details now?")

	opts := []string{
		"Enter `Key ID` and `Access Key`",
		"I don't have a `Key ID` and `Access Key` yet",
		"I don't have an AWS account yet",
	}

	switch common.Select("Select Action", opts) {
	case opts[1]:
		sc := bufio.NewScanner(os.Stdin)
		showAWSAccessKeyObtain(sc)

	case opts[2]:
		sc := bufio.NewScanner(os.Stdin)

		fmt.Println(common.AWSPrefix(), "Visit",
			common.Color("https://portal.aws.amazon.com/billing/signup#/start", "66C2CD"),
			"and create a account (follow the steps)")
		fmt.Println(common.AWSPrefix(), "< press", common.Color("enter", "#A8CC8C"), "to continue >")
		sc.Scan()

		showAWSAccessKeyObtain(sc)
	}
	return askAWSCred()
}
