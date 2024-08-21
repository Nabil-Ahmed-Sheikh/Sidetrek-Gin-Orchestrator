package emissary

type (
	CreateEmissaryMappingInput struct {
		NamespaceName string
		ClusterName   string
		ServiceName   string
		HostName      string
	}

	DestroyEmissaryMappingInput struct {
		NamespaceName string
		ClusterName   string
		ServiceName   string
		HostName      string
	}
)
