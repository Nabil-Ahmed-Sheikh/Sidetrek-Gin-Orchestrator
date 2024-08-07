package workflows

import (
	"GinProject/app/temporal/aws/deployment"
	"GinProject/app/temporal/aws/ec2"
	"GinProject/app/temporal/aws/ecr"
	"GinProject/app/temporal/aws/network"
	"GinProject/app/temporal/dagster"
	"GinProject/app/temporal/namespace"
	"encoding/json"
	"errors"

	// "go.temporal.io/api/namespace/v1"
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

	w.RegisterActivity(namespace.CreateNamespaceActivity)
	w.RegisterActivity(namespace.DestroyNamespaceActivity)

	w.RegisterActivity(ecr.ApplyDockerBuildAndPushEcrActivity)
	w.RegisterActivity(ecr.DestroyDockerBuildAndPushEcrActivity)

	w.RegisterActivity(deployment.CreateDeploymentActivity)
	w.RegisterActivity(deployment.DestroyDeploymentActivity)

	w.RegisterWorkflow(ec2.DeployEc2Workflow)
	w.RegisterWorkflow(ec2.DestroyEc2Workflow)

	w.RegisterWorkflow(network.CreateVpcWorkflow)
	w.RegisterWorkflow(network.DestroyVpcWorkflow)

	// dagster
	w.RegisterWorkflow(dagster.CreateDagsterClusterWorkflow)
	w.RegisterWorkflow(dagster.DestroyDagsterClusterWorkflow)

	// Namespace
	w.RegisterWorkflow(namespace.CreateNamespaceWorkflow)
	w.RegisterWorkflow(namespace.DestroyNamespaceWorkflow)

	// ECR

	w.RegisterWorkflow(ecr.ApplyDockerBuildAndPushEcrWorkflow)
	w.RegisterWorkflow(ecr.DestroyDockerBuildAndPushEcrWorkflow)

	w.RegisterWorkflow(deployment.CreateDeploymentWorkflow)
	w.RegisterWorkflow(deployment.DestroyDeploymentWorkflow)

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
	case "CreateNamespaceWorkflow":
		in := namespace.CreateNamespaceInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyNamespaceWorkflow":
		in := namespace.DestroyNamespaceInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "ApplyDockerBuildAndPushEcrWorkflow":
		in := ecr.ApplyDockerBuildAndPushEcrInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyDockerBuildAndPushEcrWorkflow":
		in := ecr.DestroyDockerBuildAndPushEcrInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "CreateDeploymentWorkflow":
		in := deployment.CreateDeploymentInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err
	case "DestroyDeploymentWorkflow":
		in := deployment.DestroymentInput{}
		err := json.Unmarshal([]byte(input), &in)
		return in, err

	default:
		return nil, ErrWorkflowNotFound
	}

}
