package tsess

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

const (
	DefaultSecurityGroupDescription = "Allows SSH and the Jamulus default port"
)

var (
	SecurityGroupAllowSSH = (&ec2.IpPermission{}).
				SetIpProtocol("tcp").
				SetFromPort(22).
				SetToPort(22).
				SetIpRanges([]*ec2.IpRange{
			(&ec2.IpRange{}).
				SetCidrIp("0.0.0.0/0"),
		})
	SecurityGroupAllowJamulus = (&ec2.IpPermission{}).
					SetIpProtocol("udp").
					SetFromPort(22124).
					SetToPort(22124).
					SetIpRanges([]*ec2.IpRange{
			(&ec2.IpRange{}).
				SetCidrIp("0.0.0.0/0"),
		})
)

func (s *Session) FindSecurityGroup() (*ec2.SecurityGroup, error) {
	call, err := s.EC2.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return nil, err
	}
	for _, sg := range call.SecurityGroups {
		if sg == nil || sg.GroupName == nil {
			continue
		}
		if strings.EqualFold(*sg.GroupName, s.Instance.SecurityGroupName) {
			return sg, nil
		}
	}
	return nil, errors.New("security group not found")
}

func (s *Session) CreateSecurityGroup() (err error) {
	// create security group
	var resp *ec2.CreateSecurityGroupOutput
	if resp, err = s.EC2.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description: aws.String(DefaultSecurityGroupDescription),
		GroupName:   aws.String(s.Instance.SecurityGroupName),
	}); err != nil {
		return
	}
	// add rules
	if _, err = s.EC2.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: resp.GroupId,
		IpPermissions: []*ec2.IpPermission{
			SecurityGroupAllowSSH,
			SecurityGroupAllowJamulus,
		},
	}); err != nil {
		return
	}
	return
}
