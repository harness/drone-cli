package transform

import (
	"os"
	"strings"

	"github.com/drone/drone-runtime/engine"
)

// WithProxy is a transform function that adds the
// http_proxy environment variables to every container.
func WithProxy() func(*engine.Spec) {
	environ := map[string]string{}
	if value := getenv("no_proxy"); value != "" {
		environ["no_proxy"] = value
		environ["NO_PROXY"] = value
	}
	if value := getenv("http_proxy"); value != "" {
		environ["http_proxy"] = value
		environ["HTTP_PROXY"] = value
	}
	if value := getenv("https_proxy"); value != "" {
		environ["https_proxy"] = value
		environ["HTTPS_PROXY"] = value
	}
	return WithEnviron(environ)
}

func getenv(name string) (value string) {
	name = strings.ToUpper(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	name = strings.ToLower(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	return
}
