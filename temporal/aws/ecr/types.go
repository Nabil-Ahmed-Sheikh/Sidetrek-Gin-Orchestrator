package ecr

type (
	ApplyDockerBuildAndPushEcrInput struct {
		CreateEcrRepo          bool   // create_ecr_repo
		CreateSamMetadata      bool   // create_sam_metadata
		UseImageTag            bool   // use_image_tag -> default true
		EcrAddress             string // ecr_address
		EcrUser                string // ecr_user -> default AWS
		EcrPassword            string // ecr_password -> auth token
		EcrRepoName            string // ecr_repo
		EcrImageTag            string // image_tag
		DockerFilePath         string // docker_file_path -> default Dockerfile
		SourcePath             string // source_path
		ImageTagMutability     string // image_tag_mutability -> default MUTABLE
		ScanOnPush             bool   // scan_on_push -> default false
		EcrForceDelete         bool   // ecr_force_delete -> default true
		EcrRepoTags            string // ecr_repo_tags -> default {}
		BuildArgs              string // build_args -> default {}
		EcrRepoLifecyclePolity string // ecr_repo_lifecycle_policy
		KeepRemotely           bool   // keep_remotely -> default false
		Platform               string // platform
		ForceRemove            bool   // force_remove -> default false
		KeepLocally            bool   // keep_locally -> default false
		Triggers               string // triggers -> default {}
	}

	ApplyDockerBuildAndPushEcrOutput struct {
		ImageUri string
		ImageId  string
	}

	DestroyDockerBuildAndPushEcrInput struct {
		EcrRepoName  string
		EcrImageName string
	}
)
