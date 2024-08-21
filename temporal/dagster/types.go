package dagster

type (
	CreateDagsterClusterInput struct {
		ClusterName   string
		AdditionalSet []AdditionalSet
		repository    string
	}

	AdditionalSet struct {
		Name  string
		Value string
		Type  string
	}

	DestroyDagsterClusterInput struct {
		ClusterName string
	}
)
