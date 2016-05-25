package transform

import "github.com/drone/drone-exec/yaml"

// TransformFunc defines an operation for transforming the Yaml file.
type TransformFunc func(*yaml.Config) error
