package network

import "fmt"

type (
	CreateVpcInput struct {
		Name      string
		CIDRBlock string
		Subnets   []Subnet
	}

	CreateVpcOutput struct {
		VpcID string
	}

	DestroyVpcInput struct {
		Name   string
		Region string
	}

	Subnet struct {
		AvailabilityZone string
		CIDRBlock        string
		Public           bool
	}

	Subnets []Subnet

	CreateSubnetsInput struct {
		VpcID   string
		Subnets Subnets
	}

	CreateSubnetsOutput struct{}

	DestroySubnetsInput struct {
		VpcID string
	}

	CreateSecurityGroupInput struct {
		VpcID       string
		Name        string
		Description string
	}

	DestroySecurityGroupInput struct {
		VpcID       string
		Name        string
		Description string
	}

	CreateElasticIpInput struct {
		InternetGatewayID string
	}

	DestroyElasticIpInput struct {
		InternetGatewayID string
	}

	CreateInternetGatewayInput struct {
		VpcID string
	}

	DestroyInternetGatewayInput struct {
		VpcID string
	}

	CreateNatGatewayInput struct {
		AllocationID string
		SubnetID     string
	}

	DestroyNatGatewayInput struct {
		AllocationID string
		SubnetID     string
	}

	CreateRouteTableInput struct {
		VpcID string
		Name  string
	}

	DestroyRouteTableInput struct {
		VpcID string
		Name  string
	}

	CreateRouteTableAssociationInput struct {
		SubnetID     string
		RouteTableID string
	}

	DestroyRouteTableAssociationInput struct {
		SubnetID     string
		RouteTableID string
	}

	CreateRouteInput struct {
		RouteTableID      string
		NatGatewayID      string
		InternetGatewayID string
	}

	DestroyRouteInput struct {
		RouteTableID      string
		NatGatewayID      string
		InternetGatewayID string
	}
)

func (s Subnets) String() string {
	var str string
	for _, s := range s {
		str += fmt.Sprintf("%s-%s-%t ", s.AvailabilityZone, s.CIDRBlock, s.Public)
	}
	return str
}
