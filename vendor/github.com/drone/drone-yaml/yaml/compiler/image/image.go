package image

import "github.com/docker/distribution/reference"

// Trim returns the short image name without tag.
func Trim(name string) string {
	ref, err := reference.ParseAnyReference(name)
	if err != nil {
		return name
	}
	named, err := reference.ParseNamed(ref.String())
	if err != nil {
		return name
	}
	named = reference.TrimNamed(named)
	return reference.FamiliarName(named)
}

// Expand returns the fully qualified image name.
func Expand(name string) string {
	ref, err := reference.ParseAnyReference(name)
	if err != nil {
		return name
	}
	named, err := reference.ParseNamed(ref.String())
	if err != nil {
		return name
	}
	named = reference.TagNameOnly(named)
	return named.String()
}

// Match returns true if the image name matches
// an image in the list. Note the image tag is not used
// in the matching logic.
func Match(from string, to ...string) bool {
	from = Trim(from)
	for _, match := range to {
		if from == Trim(match) {
			return true
		}
	}
	return false
}

// MatchHostname returns true if the image hostname
// matches the specified hostname.
func MatchHostname(image, hostname string) bool {
	ref, err := reference.ParseAnyReference(image)
	if err != nil {
		return false
	}
	named, err := reference.ParseNamed(ref.String())
	if err != nil {
		return false
	}
	if hostname == "index.docker.io" {
		hostname = "docker.io"
	}
	return reference.Domain(named) == hostname
}
