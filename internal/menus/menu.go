package menus

import "github.com/aws/aws-sdk-go/service/ec2"

type Menu struct {
	Parent *Menu
	Print  func()
}

func (m *Menu) Back() {
	// there is no way back
	if m.Parent == nil {
		return
	}
	m.Parent.Print()
}

// EC2

type EC2Menu struct {
	*Menu
	ec *ec2.EC2
}
