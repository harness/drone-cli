package build

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"strconv"
	"strings"
)

// parses a string into a build #, with no support for resolving the latest build
// returns very friendly errors as opposed to strconv errors
func parseBuildArg(arg string) (buildNumber int, err error) {
	if len(arg) == 0 {
		err = fmt.Errorf("Error: missing build number after repository.")
		return
	}
	buildNumber, err = strconv.Atoi(arg)
	if err != nil {
		err = fmt.Errorf("Error: malformed build number specifed: %s", err.Error())
	}

	return
}

// parses a string, which is one of "last" or an exact build number.
// the returned value is the latest build, it's number
//
// this is redundant but more often than not you don't actually want/need the actual build,
// just it's number; you can use _ to throw away the build and 'number' to access just the number
//
// this is an example
//
// _, number, err := getBuildWithArg(c.Args().Get(1), owner, repo, client)
//
func getBuildWithArg(arg, owner, repo string, client drone.Client) (build *drone.Build, number int, err error) {
	arg = strings.ToLower(arg)
	switch arg {
	case "last":
		build, err = client.BuildLast(owner, repo, "")
		if err != nil {
			err = fmt.Errorf("Error: could not get last build: %s", err.Error())
		}
		number = build.Number
	default:
		number, err = parseBuildArg(arg)
		if err != nil {
			return
		}
		build, err = client.Build(owner, repo, number)
		if err != nil {
			err = fmt.Errorf("Error: Could not get build %d: %s", number, err.Error())
		}
	}
	return
}
