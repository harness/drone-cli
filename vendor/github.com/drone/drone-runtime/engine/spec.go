// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package engine

type (
	// Metadata provides execution metadata.
	Metadata struct {
		UID       string            `json:"uid,omitempty"`
		Namespace string            `json:"namespace,omitempty"`
		Name      string            `json:"name,omitempty"`
		Labels    map[string]string `json:"labels,omitempty"`
	}

	// Spec provides the pipeline spec. This provides the
	// required instructions for reproducable pipeline
	// execution.
	Spec struct {
		Metadata Metadata  `json:"metadata,omitempty"`
		Platform Platform  `json:"platform,omitempty"`
		Secrets  []*Secret `json:"secrets,omitempty"`
		Steps    []*Step   `json:"steps,omitempty"`
		Files    []*File   `json:"files,omitempty"`

		// Docker-specific settings. These settings are
		// only used by the Docker and Kubernetes runtime
		// drivers.
		Docker *DockerConfig `json:"docker,omitempty"`

		// Qemu-specific settings. These settings are only
		// used by the qemu runtime driver.
		Qemu *QemuConfig `json:"qemu,omitempty"`

		// VMWare Fusion settings. These settings are only
		// used by the VMWare runtime driver.
		Fusion *FusionConfig `json:"fusion,omitempty"`
	}

	// Step defines a pipeline step.
	Step struct {
		Metadata     Metadata          `json:"metadata,omitempty"`
		Detach       bool              `json:"detach,omitempty"`
		DependsOn    []string          `json:"depends_on,omitempty"`
		Devices      []*VolumeDevice   `json:"devices,omitempty"`
		Envs         map[string]string `json:"environment,omitempty"`
		Files        []*FileMount      `json:"files,omitempty"`
		IgnoreErr    bool              `json:"ignore_err,omitempty"`
		IgnoreStdout bool              `json:"ignore_stderr,omitempty"`
		IgnoreStderr bool              `json:"ignore_stdout,omitempty"`
		Resources    *Resources        `json:"resources,omitempty"`
		RunPolicy    RunPolicy         `json:"run_policy,omitempty"`
		Secrets      []*SecretVar      `json:"secrets,omitempty"`
		Volumes      []*VolumeMount    `json:"volumes,omitempty"`
		WorkingDir   string            `json:"working_dir,omitempty"`

		// Docker-specific settings. These settings are
		// only used by the Docker and Kubernetes runtime
		// drivers.
		Docker *DockerStep `json:"docker,omitempty"`
	}

	// DockerAuth defines dockerhub authentication credentials.
	DockerAuth struct {
		Address  string `json:"address,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	// DockerConfig configures a Docker-based pipeline.
	DockerConfig struct {
		Auths   []*DockerAuth `json:"auths,omitempty"`
		Volumes []*Volume     `json:"volumes,omitempty"`
	}

	// DockerStep configures a docker step.
	DockerStep struct {
		Args       []string   `json:"args,omitempty"`
		Command    []string   `json:"command,omitempty"`
		DNS        []string   `json:"dns,omitempty"`
		DNSSearch  []string   `json:"dns_search,omitempty"`
		ExtraHosts []string   `json:"extra_hosts,omitempty"`
		Image      string     `json:"image,omitempty"`
		Networks   []string   `json:"networks,omitempty"`
		Ports      []*Port    `json:"ports,omitempty"`
		Privileged bool       `json:"privileged,omitempty"`
		PullPolicy PullPolicy `json:"pull_policy,omitempty"`
		User       string     `json:"user"`
	}

	// File defines a file that should be uploaded or
	// mounted somewhere in the step container or virtual
	// machine prior to command execution.
	File struct {
		Metadata Metadata `json:"metadata,omitempty"`
		Data     []byte   `json:"data,omitempty"`
	}

	// FileMount defines how a file resource should be
	// mounted or included in the runtime environment.
	FileMount struct {
		Name string `json:"name,omitempty"`
		Path string `json:"path,omitempty"`
		Mode int64  `json:"mode,omitempty"`

		// Base string `json:"base,omitempty"`
	}

	// FusionConfig configures a VMWare Fusion-based pipeline.
	FusionConfig struct {
		Image string `json:"image,omitempty"`
	}

	// Platform defines the target platform.
	Platform struct {
		OS      string `json:"os,omitempty"`
		Arch    string `json:"arch,omitempty"`
		Variant string `json:"variant,omitempty"`
		Version string `json:"version,omitempty"`
	}

	// Port represents a network port in a single container.
	Port struct {
		Port     int    `json:"port,omitempty"`
		Host     int    `json:"host,omitempty"`
		Protocol string `json:"protocol,omitempty"`
	}

	// QemuConfig configures a Qemu-based pipeline.
	QemuConfig struct {
		Image string `json:"image,omitempty"`
	}

	// Resources describes the compute resource
	// requirements.
	Resources struct {
		// Limits describes the maximum amount of compute
		// resources allowed.
		Limits *ResourceObject `json:"limits,omitempty"`

		// Requests describes the minimum amount of
		// compute resources required.
		Requests *ResourceObject `json:"requests,omitempty"`
	}

	// ResourceObject describes compute resource
	// requirements.
	ResourceObject struct {
		CPU    int64 `json:"cpu,omitempty"`
		Memory int64 `json:"memory,omitempty"`
	}

	// Secret represents a secret variable.
	Secret struct {
		Metadata Metadata `json:"metadata,omitempty"`
		Data     string   `json:"data,omitempty"`
	}

	// SecretVar represents an environment variable
	// sources from a secret.
	SecretVar struct {
		Name string `json:"name,omitempty"`
		Env  string `json:"env,omitempty"`
	}

	// State represents the container state.
	State struct {
		ExitCode  int  // Container exit code
		Exited    bool // Container exited
		OOMKilled bool // Container is oom killed
	}

	// Volume that can be mounted by containers.
	Volume struct {
		Metadata Metadata        `json:"metadata,omitempty"`
		EmptyDir *VolumeEmptyDir `json:"temp,omitempty"`
		HostPath *VolumeHostPath `json:"host,omitempty"`
	}

	// VolumeDevice describes a mapping of a raw block
	// device within a container.
	VolumeDevice struct {
		Name       string `json:"name,omitempty"`
		DevicePath string `json:"path,omitempty"`
	}

	// VolumeMount describes a mounting of a Volume
	// within a container.
	VolumeMount struct {
		Name string `json:"name,omitempty"`
		Path string `json:"path,omitempty"`
	}

	// VolumeEmptyDir mounts a temporary directory from the
	// host node's filesystem into the container. This can
	// be used as a shared scratch space.
	VolumeEmptyDir struct {
		Medium    string `json:"medium,omitempty"`
		SizeLimit int64  `json:"size_limit,omitempty"`
	}

	// VolumeHostPath mounts a file or directory from the
	// host node's filesystem into your container.
	VolumeHostPath struct {
		Path string `json:"path,omitempty"`
	}
)
