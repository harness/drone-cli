package builder

import (
	"fmt"
	"io"

	"github.com/drone/drone-cli/common/uuid"
)

const (
	ImageInit  = "drone/drone-init"
	ImageClone = "drone/drone-clone-git"
)

func Run(req *Request, resp ResponseWriter) error {

	var containers []*Container

	defer func() {
		for i := len(containers) - 1; i >= 0; i-- {
			container := containers[i]
			container.Stop()
			container.Kill()
			container.Remove()
		}
	}()

	// temporary name for the build container
	//name := fmt.Sprintf("build-init-%s", createUID())
	net := req.Config.Docker.Net
	uid := uuid.CreateUUID()
	cmd := []string{req.Encode()}

	// init container
	containers = append(containers, &Container{
		Name:    fmt.Sprintf("drone-%s-init", uid),
		Image:   ImageInit,
		Volumes: []string{"/drone"},
		Cmd:     cmd,
	})

	// clone container
	containers = append(containers, &Container{
		Name:        fmt.Sprintf("drone-%s-clone", uid),
		Image:       ImageClone,
		VolumesFrom: []string{containers[0].Name},
		Cmd:         cmd,
	})

	// attached service containers
	for i, service := range req.Config.Services {
		containers = append(containers, &Container{
			Name:        fmt.Sprintf("drone-%s-service-%v", uid, i),
			Image:       service,
			Env:         req.Config.Env,
			NetworkMode: net,
			Detached:    true,
		})

		if i == 0 && len(net) == 0 {
			net = fmt.Sprintf("container:drone-%s-service-%v", uid, i)
		}
	}

	// build container
	containers = append(containers, &Container{
		Name:        fmt.Sprintf("drone-%s-build", uid),
		Image:       req.Config.Image,
		Env:         req.Config.Env,
		Cmd:         []string{"/drone/bin/build.sh"},
		Entrypoint:  []string{"/bin/bash"},
		WorkingDir:  req.Clone.Dir,
		NetworkMode: net,
		Privileged:  req.Config.Docker.Privileged,
		VolumesFrom: []string{containers[0].Name},
	})

	//
	// create the notify, publish, deploy containers
	//

	// loop through and create containers
	for _, container := range containers {
		container.SetClient(req.Client)
		if err := container.Create(); err != nil {
			return err
		}
	}

	// loop through and start containers
	for _, container := range containers {
		if err := container.Start(); err != nil {
			return err
		}
		if container.Detached { // if a detached (daemon) just continue
			continue
		}
		r, err := container.Logs()
		if err != nil {
			return err
		}
		io.Copy(resp, r)
		r.Close()
		info, err := container.Inspect()
		if err != nil {
			return err
		}

		if info.State.Running != false {
			fmt.Println("ERROR: container still running")
		}

		resp.WriteExitCode(info.State.ExitCode)
		if info.State.ExitCode != 0 {
			break
		}
	}

	return nil
}
