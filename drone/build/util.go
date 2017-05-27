package build

import (
	"fmt"
	"strconv"
	"github.com/drone/drone-go/drone"
	"strings"
)

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