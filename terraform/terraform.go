package terraform

import "embed"

//go:embed aws
var AWS embed.FS

//go:embed dagster
var Dagster embed.FS

//go:embed namespace
var Namespace embed.FS
