# drone-cli
Drone command-line interface

### System Requirements

The Drone command-line utility requires **Docker 1.6** or higher. If you are using Windows or Mac you can install Docker using the [Docker Toolbox](https://www.docker.com/docker-toolbox).

### Installation

**Linux**

Download and install the x64 linux binary:

```
curl http://downloads.drone.io/drone-cli/drone_linux_amd64.tar.gz | tar zx
sudo install -t /usr/local/bin drone
```

**OSX**

Download and install using Homebrew:

```
brew tap drone/drone
brew install drone
```

Or manually download and install the binary:

```
curl http://downloads.drone.io/drone-cli/drone_darwin_amd64.tar.gz | tar zx
sudo cp drone /usr/local/bin
```

### Authentication

You must provide the command line utility with the Drone server URL and a valid API token. You can get your API token in the Drone profile page.

```
export DRONE_TOKEN=
export DRONE_SERVER=
```

