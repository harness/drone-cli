package transform

import (
	"path/filepath"

	"github.com/drone/drone-exec/yaml"
	"github.com/drone/drone-go/drone"
)

func ImageSecrets(c *yaml.Config, secrets []*drone.Secret, event string) error {
	var images []*yaml.Container
	images = append(images, c.Pipeline...)
	images = append(images, c.Services...)

	for _, image := range images {
		imageSecrets(image, secrets, event)
	}
	return nil
}

func imageSecrets(c *yaml.Container, secrets []*drone.Secret, event string) {
	for _, secret := range secrets {
		if !match(secret, c.Image, event) {
			continue
		}

		switch secret.Name {
		case "REGISTRY_USERNAME":
			c.AuthConfig.Username = secret.Value
		case "REGISTRY_PASSWORD":
			c.AuthConfig.Password = secret.Value
		case "REGISTRY_EMAIL":
			c.AuthConfig.Email = secret.Value
		default:
			if c.Environment == nil {
				c.Environment = map[string]string{}
			}
			c.Environment[secret.Name] = secret.Value
		}
	}
}

// match returns true if an image and event match the restricted list.
func match(s *drone.Secret, image, event string) bool {
	return matchImage(s, image) && matchEvent(s, event)
}

// matchImage returns true if an image matches the restricted list.
func matchImage(s *drone.Secret, image string) bool {
	for _, pattern := range s.Images {
		if match, _ := filepath.Match(pattern, image); match {
			return true
		} else if pattern == "*" {
			return true
		}
	}
	return false
}

// matchEvent returns true if an event matches the restricted list.
func matchEvent(s *drone.Secret, event string) bool {
	for _, pattern := range s.Events {
		if match, _ := filepath.Match(pattern, event); match {
			return true
		}
	}
	return false
}
