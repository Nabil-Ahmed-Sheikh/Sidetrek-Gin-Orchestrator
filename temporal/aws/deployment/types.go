package deployment

type (
	CreateDeploymentInput struct {
		EcrAddress  string
		ClusterName string
	}

	DestroymentInput struct {
		EcrAddress string
	}
)
