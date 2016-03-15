# drone-go

[![Build Status](http://beta.drone.io/api/badges/drone/drone-go/status.svg)](http://beta.drone.io/drone/drone-go)

drone-go is a Go client library for accessing the Drone [API](http://readme.drone.io/devs/api/builds/) and writing [plugins](http://readme.drone.io/plugins/).

Download the package using `go get`:

```Go
go get "github.com/drone/drone-go/drone"
go get "github.com/drone/drone-go/plugin"
```

Import the package:

```Go
import "github.com/drone/drone-go/drone"
```

Create the client:

```Go
const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	host  = "http://drone.company.com"
)

client := drone.NewClientToken(host, token)
```

Get the current user:

```Go
user, err := client.Self()
fmt.Println(user)
```

Get the repository list:

```Go
repos, err := client.RepoList()
fmt.Println(repos)
```

Get the named repository:

```Go
repo, err := client.Repo("drone", "drone-go")
fmt.Println(repo)
```
