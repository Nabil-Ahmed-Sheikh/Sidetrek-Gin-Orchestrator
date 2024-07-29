package ec2

import "fmt"

type (
	CreateEc2Input struct {
	}

	CreateEc2Output struct {
	}

	Subnet struct {
		AvailabilityZone string
		CIDRBlock        string
		Public           bool
	}

	DestroyEc2Input struct {
	}

	DestroyEc2Output struct {
	}

	Subnets []Subnet
)

func (s Subnets) String() string {
	var str string
	for _, s := range s {
		str += fmt.Sprintf("%s-%s-%t ", s.AvailabilityZone, s.CIDRBlock, s.Public)
	}
	return str
}
