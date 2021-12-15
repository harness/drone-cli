module github.com/drone/drone-cli

go 1.16

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/docker/go-units v0.4.0
	github.com/drone/drone-go v1.7.0
	github.com/drone/drone-runtime v1.1.1-0.20200623162453-61e33e2cab5d
	github.com/drone/drone-yaml v0.0.0-20190729072335-70fa398b3560
	github.com/drone/envsubst v1.0.3
	github.com/drone/funcmap v0.0.0-20190918184546-d4ef6e88376d
	github.com/drone/signal v1.0.0
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-jsonnet v0.17.0
	github.com/jackspirou/syscerts v0.0.0-20160531025014-b68f5469dff1
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-colorable v0.1.4
	github.com/mattn/go-isatty v0.0.11
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.2
	go.starlark.net v0.0.0-20201118183435-e55f603d8c79
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
