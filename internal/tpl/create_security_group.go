package tpl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
	"log"
	"time"
)

const DefaultSecurityGroup = "jamulus-security-group"

func (t *CreateInstanceTemplate) CreateSecurityGroup(ec *ec2.EC2) {
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Prefix = "🤔 Creating security group "
	s.FinalMSG = "😁 Created security group!"
	s.Start()
	defer s.Stop()

	// create security group
	createResponse, err := ec.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description:       aws.String("Allows SSH and the Jamulus default port"),
		GroupName:         aws.String(DefaultSecurityGroup),
		TagSpecifications: nil,
		VpcId:             nil,
	})
	if err != nil {
		log.Fatalln("Error creating security group:", err)
		return
	}

	// update security group
	t.Instance.SecurityGroupID = *createResponse.GroupId

	// assign rules to security group
	if _, err := ec.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: &t.Instance.SecurityGroupID,
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{}).
				SetIpProtocol("tcp").
				SetFromPort(22).
				SetToPort(22).
				SetIpRanges([]*ec2.IpRange{
					(&ec2.IpRange{}).
						SetCidrIp("0.0.0.0/0"),
				}),
			(&ec2.IpPermission{}).
				SetIpProtocol("udp").
				SetFromPort(22124).
				SetToPort(22124).
				SetIpRanges([]*ec2.IpRange{
					(&ec2.IpRange{}).
						SetCidrIp("0.0.0.0/0"),
				}),
		},
	}); err != nil {
		log.Fatalln("Error assigning rules to security group:", err)
		return
	}
}
