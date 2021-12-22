module github.com/drone/drone-cli

go 1.17

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/buildkite/yaml v2.1.0+incompatible
	github.com/docker/go-units v0.4.0
	github.com/drone-runners/drone-runner-docker v1.8.0
	github.com/drone/drone-go v1.7.1
	github.com/drone/drone-runtime v1.1.1-0.20200623162453-61e33e2cab5d
	github.com/drone/drone-yaml v0.0.0-20190729072335-70fa398b3560
	github.com/drone/envsubst v1.0.3
	github.com/drone/funcmap v0.0.0-20190918184546-d4ef6e88376d
	github.com/drone/runner-go v1.12.0
	github.com/drone/signal v1.0.0
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-jsonnet v0.17.0
	github.com/jackspirou/syscerts v0.0.0-20160531025014-b68f5469dff1
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-colorable v0.1.4
	github.com/mattn/go-isatty v0.0.11
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.2
	go.starlark.net v0.0.0-20201118183435-e55f603d8c79
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

require (
	github.com/99designs/httpsignatures-go v0.0.0-20170731043157-88528bf4ca7e // indirect
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/bmatcuk/doublestar v1.1.1 // indirect
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/natessilva/dag v0.0.0-20180124060714-7194b8dcc5c4 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20201110150050-8816d57aaa9a // indirect
	google.golang.org/grpc v1.33.2 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
