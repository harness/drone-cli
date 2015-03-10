package runner

/*
func Notify(req *Request, resp ResponseWriter) error {

	var containers []*Container

	defer func() {
		for i := len(containers) - 1; i >= 0; i-- {
			container := containers[i]
			container.Stop()
			container.Kill()
			container.Remove()
		}
	}()

	// attached service containers
	for _, notification := range req.Config.Notify {
		containers = append(containers, &Container{
			Image: notification.Image,
			Env:   notification.Environment,
			Cmd:   EncodeParams(req, notification.Config),
		})
	}

	// loop through and run containers
	for _, container := range containers {
		container.SetClient(req.Client)
		if err := container.Create(); err != nil {
			return err
		}
		if err := container.Start(); err != nil {
			return err
		}
		container.Wait()
	}

	return nil
}
*/
