package workflows

import (
	"GinProject/app/temporal/aws/ec2"
	"GinProject/app/temporal/aws/network"
	"GinProject/app/temporal/dagster"
	"encoding/json"
	"errors"

	"go.temporal.io/sdk/worker"
)

var (
	ErrWorkflowNotFound = errors.New("workflow not found")
)

func Register(w worker.Worker) {
	w.RegisterActivity(ec2.CreateEc2Activity)
	w.RegisterActivity(ec2.DestroyEc2Activity)

	w.RegisterActivity(network.CreateVpcActivity)
	w.RegisterActivity(network.DestroyVpcActivity)

	w.RegisterActivity(dagster.CreateDagsterClusterActivity)
	w.RegisterActivity(dagster.DestroyDagsterClusterActivity)

	w.RegisterWorkflow(ec2.DeployEc2Workflow)
	w.RegisterWorkflow(ec2.DestroyEc2Workflow)

	w.RegisterWorkflow(network.CreateVpcWorkflow)
	w.RegisterWorkflow(network.DestroyVpcWorkflow)

	// dagster
	w.RegisterWorkflow(dagster.CreateDagsterClusterWorkflow)
	w.RegisterWorkflow(dagster.DestroyDagsterClusterWorkflow)
}

func MapInputToWorkflow(wfname string, input string) (interface{}, error) {
	switch wfname {
	case "CreateVpcWorkflow":
		in := network.CreateVpcInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyVpcWorkflow":
		in := network.DestroyVpcInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DeployEc2Workflow":
		in := ec2.CreateEc2Input{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyEc2Workflow":
		in := ec2.DestroyEc2Output{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "CreateDagsterClusterWorkflow":
		in := dagster.CreateDagsterClusterInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyDagsterClusterWorkflow":
		in := dagster.DestroyDagsterClusterInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	default:
		return nil, ErrWorkflowNotFound
	}

}
