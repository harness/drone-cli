package compiler

import (
	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"

	"github.com/dchest/uniuri"
)

// TODO(bradrydzewski) handle depends_on (clone, services, etc)

// random provides the default function used to
// generate a random string.
var random = uniuri.New

// A Compiler compiles the pipeline configuration to an
// intermediate representation that can be executed by
// the Drone runtime engine.
type Compiler struct {
	// GitCredentialsFunc returns a .git-credentials file
	// that can be used by the default clone step to
	// authenticate to the remote repository.
	GitCredentialsFunc func() []byte

	// NetrcFunc returns a .netrc file that can be used by
	// the default clone step to authenticate to the remote
	// repository.
	NetrcFunc func() []byte

	// PrivilegedFunc returns true if the container should
	// be started in privileged mode. The intended use is
	// for plugins that run Docker-in-Docker. This will be
	// deprecated in a future release.
	PrivilegedFunc func(*yaml.Container) bool

	// RandFunc returns a random string. The random
	// function is used to create unique identifiers for
	// the namespace, container, and volume resources.
	RandFunc func() string

	// SkipFunc returns true if the step should be skipped.
	// The skip function can be used to evaluate the when
	// clause of each step, and return true if it should
	// be skipped.
	SkipFunc func(*yaml.Container) bool

	// TransformFunc can be used to modify the compiled
	// output prior to completion. This can be useful when
	// you need to programatically modify the output,
	// set defaults, etc.
	TransformFunc func(spec *engine.Spec)
}

// Compile returns an intermediate representation of the
// pipeline configuration that can be executed by the
// Drone runtime engine.
func (c *Compiler) Compile(from *yaml.Pipeline) *engine.Spec {
	namespace := c.random()

	spec := &engine.Spec{
		Metadata: engine.Metadata{
			UID:       namespace,
			Name:      namespace,
			Namespace: namespace,
			Labels: map[string]string{
				"io.drone.pipeline.name": from.Name,
				"io.drone.pipeline.kind": from.Kind,
				"io.drone.pipeline.type": from.Type,
			},
		},
		Platform: engine.Platform{
			OS:      from.Platform.OS,
			Arch:    from.Platform.Arch,
			Version: from.Platform.Version,
			Variant: from.Platform.Variant,
		},
		Docker:  &engine.DockerConfig{},
		Files:   nil,
		Secrets: nil,
	}

	// create the default workspace path. If a container
	// does not specify a working directory it defaults
	// to the workspace path.
	base, dir, workspace := createWorkspace(from)

	// create the default workspace volume definition.
	// the volume will be mounted to each container in
	// the pipeline.
	spec.Docker.Volumes = append(spec.Docker.Volumes,
		&engine.Volume{
			Metadata: engine.Metadata{
				UID:       c.random(),
				Name:      workspaceName,
				Namespace: namespace,
				Labels:    map[string]string{},
			},
			EmptyDir: &engine.VolumeEmptyDir{},
		},
	)

	// for each volume defined in the yaml configuration
	// file, convert to a runtime volume and append to the
	// specification.
	for _, from := range from.Volumes {
		to := &engine.Volume{
			Metadata: engine.Metadata{
				UID:       c.random(),
				Name:      from.Name,
				Namespace: namespace,
				Labels:    map[string]string{},
			},
		}
		if from.EmptyDir != nil {
			// if the yaml configuration specifies an empty
			// directory volume (data volume) or an in-memory
			// file system.
			to.EmptyDir = &engine.VolumeEmptyDir{
				Medium:    from.EmptyDir.Medium,
				SizeLimit: int64(from.EmptyDir.SizeLimit),
			}
		} else if from.HostPath != nil {
			// if the yaml configuration specifies a bind
			// mount to the host machine.
			to.HostPath = &engine.VolumeHostPath{
				Path: from.HostPath.Path,
			}
		}
		spec.Docker.Volumes = append(spec.Docker.Volumes, to)
	}

	if !from.Clone.Disable {
		src := createClone(from)
		dst := createStep(spec, src)
		dst.Docker.PullPolicy = engine.PullIfNotExists
		setupCloneDepth(from, dst)
		setupCloneCredentials(spec, dst, c.gitCredentials())
		setupWorkingDir(src, dst, workspace)
		setupWorkingDirMount(dst, base)
		setupWorkspaceEnv(dst, base, dir, workspace)
		spec.Steps = append(spec.Steps, dst)
	}

	// for each pipeline service defined in the yaml
	// configuration file, convert to a runtime step
	// and append to the specification.
	for _, service := range from.Services {
		step := createStep(spec, service)
		// note that all services are automatically
		// set to run in detached mode.
		step.Detach = true
		setupWorkingDir(service, step, workspace)
		setupWorkingDirMount(step, base)
		setupWorkspaceEnv(step, base, dir, workspace)
		// if the skip callback function returns true,
		// modify the runtime step to never execute.
		if c.skip(service) {
			step.RunPolicy = engine.RunNever
		}
		// if the step is a plugin and should be executed
		// in privileged mode, set the privileged flag.
		if c.privileged(service) {
			step.Docker.Privileged = true
		}
		spec.Steps = append(spec.Steps, step)
	}

	// rename will store a list of container names
	// that should be mapped to their temporary alias.
	rename := map[string]string{}

	// for each pipeline step defined in the yaml
	// configuration file, convert to a runtime step
	// and append to the specification.
	for _, container := range from.Steps {
		var step *engine.Step
		switch {
		case container.Build != nil:
			step = createBuildStep(spec, container)
			rename[container.Build.Image] = step.Metadata.UID
		default:
			step = createStep(spec, container)
		}
		setupWorkingDir(container, step, workspace)
		setupWorkingDirMount(step, base)
		setupWorkspaceEnv(step, base, dir, workspace)
		// if the skip callback function returns true,
		// modify the runtime step to never execute.
		if c.skip(container) {
			step.RunPolicy = engine.RunNever
		}
		// if the step is a plugin and should be executed
		// in privileged mode, set the privileged flag.
		if c.privileged(container) {
			step.Docker.Privileged = true
		}
		spec.Steps = append(spec.Steps, step)
	}

	// if the pipeline includes any build and publish
	// steps we should create an entry for the host
	// machine docker socket.
	if spec.Docker != nil && len(rename) > 0 {
		v := &engine.Volume{
			Metadata: engine.Metadata{
				UID:       c.random(),
				Name:      "_docker_socket",
				Namespace: namespace,
				Labels:    map[string]string{},
			},
			HostPath: &engine.VolumeHostPath{
				Path: "/var/run/docker.sock",
			},
		}
		spec.Docker.Volumes = append(spec.Docker.Volumes, v)
	}

	// images created during the pipeline are assigned a
	// random alias. All references to the origin image
	// name must be changed to the alias.
	for _, step := range spec.Steps {
		for k, v := range rename {
			if image.MatchTag(step.Docker.Image, k) {
				img := image.Trim(step.Docker.Image) + ":" + v
				step.Docker.Image = image.Expand(img)
			}
		}
	}

	// executes user-defined transformations before the
	// final specification is returned.
	if c.TransformFunc != nil {
		c.TransformFunc(spec)
	}

	return spec
}

// return a .git-credentials file. If the user-defined
// function is nil, a nil credentials file is returned.
func (c *Compiler) gitCredentials() []byte {
	if c.GitCredentialsFunc != nil {
		return c.GitCredentialsFunc()
	}
	return nil
}

// return a .netrc file. If the user-defined function is
// nil, a nil netrc file is returned.
func (c *Compiler) netrc() []byte {
	if c.NetrcFunc != nil {
		return c.NetrcFunc()
	}
	return nil
}

// return true if the step should be executed in privileged
// mode. If the user-defined privileged function is nil,
// a default value of false is returned.
func (c *Compiler) privileged(container *yaml.Container) bool {
	if c.PrivilegedFunc != nil {
		return c.PrivilegedFunc(container)
	}
	return false
}

// return true if the step should be skipped. If the
// user-defined skip function is nil, a defalt skip
// function is used that always returns true (i.e. do not skip).
func (c *Compiler) skip(container *yaml.Container) bool {
	if c.SkipFunc != nil {
		return c.SkipFunc(container)
	}
	return false
}

// return a random string. If the user-defined random
// function is nil, a defalt random function is returned.
func (c *Compiler) random() string {
	if c.RandFunc != nil {
		return c.RandFunc()
	}
	return uniuri.New()
}
