package runner

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/uuid"
)

const (
	ImageInit  = "plugins/drone-build"
	ImageClone = "plugins/drone-git"
)

// TODO(brydzews) use correct NetworkMode instead of hard-coded
// TODO(brydzews) use privileged when specified for all containers

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
	net := ""
	uid := uuid.CreateUUID()

	// init container
	containers = append(containers, &Container{
		Name:    fmt.Sprintf("drone-%s-init", uid),
		Image:   ImageInit,
		Volumes: []string{"/drone"},
		Cmd:     EncodeParams(req, req.Config.Build.Config),
	})

	// clone container
	containers = append(containers, &Container{
		Name:        fmt.Sprintf("drone-%s-clone", uid),
		Image:       ImageClone,
		VolumesFrom: []string{containers[0].Name},
		Cmd:         EncodeParams(req, nil),
	})

	// attached service containers
	i := 0
	for _, service := range req.Config.Compose {

		containers = append(containers, &Container{
			Name:        fmt.Sprintf("drone-%s-service-%v", uid, i),
			Image:       service.Image,
			Env:         service.Environment,
			Privileged:  service.Privileged,
			NetworkMode: net, //service.NetworkMode,
			Detached:    true,
		})

		if i == 0 && len(net) == 0 {
			net = fmt.Sprintf("container:drone-%s-service-%v", uid, i)
		}
		i++
	}

	// build container
	containers = append(containers, &Container{
		Name:        fmt.Sprintf("drone-%s-build", uid),
		Image:       req.Config.Build.Image,
		Env:         req.Config.Build.Environment,
		Cmd:         []string{"/drone/bin/build.sh"},
		Entrypoint:  []string{"/bin/bash"},
		WorkingDir:  req.Clone.Dir,
		NetworkMode: net,
		Privileged:  req.Config.Build.Privileged,
		VolumesFrom: []string{containers[0].Name},
	})

	//
	// TODO: create the notify, publish, deploy containers
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

func EncodeParams(req *Request, args map[string]interface{}) []string {
	out := struct {
		Clone  *common.Clone          `json:"clone"`
		Commit *common.Commit         `json:"commit"`
		Repo   *common.Repo           `json:"repo"`
		User   *common.User           `json:"user"`
		Vargs  map[string]interface{} `json:"vargs"`
	}{req.Clone, req.Commit, req.Repo, req.User, args}
	raw, _ := json.Marshal(out)
	return []string{string(raw)}
}
