package tpl

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

func GetPrettyState(state *ec2.InstanceState) string {
	switch *state.Name {
	case "pending":
		return "⏱ | pending"
	case "running":
		return "✅ | running"
	case "shutting-down":
		return "🔻 | shutting down"
	case "terminated":
		return "🗑 | terminated"
	case "stopping":
		return "🥱 | stopping"
	case "stopped":
		return "😴 | stopped"
	}
	return ""
}

func (t *CreateInstanceTemplate) OpenSession(ec *ec2.EC2, instance *ec2.Instance) (client *ssh.Client, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(t.Instance.KeyPairPath); err != nil {
		log.Println("Error reading key file:", err)
		return
	}

	var signer ssh.Signer
	if signer, err = ssh.ParsePrivateKey(data); err != nil {
		log.Println("Error parsing key:", err)
		return
	}

	config := &ssh.ClientConfig{
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// wait until instance is running
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Prefix = "🤔 Waiting for instance to be ready "
	s.FinalMSG = "😁 Instance is running!"
	s.Start()

	var host string
	for {
		resp, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{instance.InstanceId},
		})
		if err != nil {
			log.Fatalln("Error reading instance:", err)
			return nil, err
		}
		i := resp.Reservations[0].Instances[0]
		s.Prefix = fmt.Sprintf("🤔 Waiting for instance to be ready [%s] ", GetPrettyState(i.State))
		if *i.State.Name == ec2.InstanceStateNameRunning {
			host = *i.PublicDnsName
			break
		}
		time.Sleep(2 * time.Second)
	}

	s.Stop()
	fmt.Println()

	// empty host?
	if host == "" {
		log.Fatalln("Error reading server hostname! (empty)")
		return nil, nil
	}

	s.Prefix = "🤔 Waiting for SSH connection "
	s.FinalMSG = "😁 Connected to SSH"
	s.Start()

	for {
		if client, err = ssh.Dial("tcp", host+":22", config); err != nil {
			s.Prefix = "🤨 Waiting for SSH connection [" + err.Error() + "] "
			time.Sleep(time.Second)
			continue
		}
		break
	}

	s.Stop()
	fmt.Println()

	// start session
	return
}
