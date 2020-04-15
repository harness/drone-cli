module github.com/drone/drone-cli

go 1.12

require (
	docker.io/go-docker v1.0.0
	github.com/99designs/httpsignatures-go v0.0.0-20170731043157-88528bf4ca7e
	github.com/Microsoft/go-winio v0.4.11
	github.com/bmatcuk/doublestar v1.1.1
	github.com/buildkite/yaml v2.1.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go-connections v0.3.0
	github.com/docker/go-units v0.3.3
	github.com/drone/drone-go v1.3.0
	github.com/drone/drone-runtime v1.0.7-0.20190729070836-38f28a11afe8
	github.com/drone/drone-yaml v0.0.0-20190729072335-70fa398b3560
	github.com/drone/envsubst v1.0.1
	github.com/drone/funcmap v0.0.0-20190918184546-d4ef6e88376d
	github.com/drone/signal v1.0.0
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v0.0.0-20170307180453-100ba4e88506
	github.com/golang/protobuf v1.2.0
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/go-jsonnet v0.11.2
	github.com/jackspirou/syscerts v0.0.0-20160531025014-b68f5469dff1
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.4
	github.com/natessilva/dag v0.0.0-20180124060714-7194b8dcc5c4
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.8.0
	github.com/urfave/cli v1.20.0
	github.com/vinzenz/yaml v0.0.0-20170920082545-91409cdd725d
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5
	golang.org/x/net v0.0.0-20181005035420-146acd28ed58
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys v0.0.0-20191002063906-3421d5a6bb1c
	google.golang.org/appengine v1.3.0
	gopkg.in/yaml.v2 v2.2.2
)
