package namespace

type (
	CreateNamespaceInput struct {
		NamespaceName string
		ClusterName   string
	}

	DestroyNamespaceInput struct {
		NamespaceName string
		ClusterName   string
	}
)
