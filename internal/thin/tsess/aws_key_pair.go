package tsess

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"strings"
)

const KeyPairPermMode = 0644

func (s *Session) FindKeyPair() (*ec2.KeyPairInfo, error) {
	call, err := s.EC2.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, err
	}
	for _, kp := range call.KeyPairs {
		if kp == nil || kp.KeyName == nil {
			continue
		}
		if strings.EqualFold(*kp.KeyName, s.Instance.KeyPair.Name) {
			return kp, nil
		}
	}
	return nil, errors.New("key pair not found")
}

func (s *Session) CreateKeyPair() (err error) {
	var resp *ec2.CreateKeyPairOutput
	if resp, err = s.EC2.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(s.Instance.KeyPair.Name),
	}); err != nil {
		return
	}

	// find available file name
	// key.pem.1
	// key.pem.2
	// ...
	// key.pem.n
	var outPath string
	var num uint64
	for {
		outPath = fmt.Sprintf("%s.%d", s.TemplatePath, num)
		if i, e := os.Stat(outPath); os.IsNotExist(e) && !i.IsDir() {
			break
		}
		num++
	}

	// write key
	err = os.WriteFile(outPath, []byte(*resp.KeyMaterial), KeyPairPermMode)
	return
}
