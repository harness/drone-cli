package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "time"
)

func convertWindowsPath(str string) string {
	// Split path into Array
	var pwds = strings.Split(str, "\\")

	// Convert drive "C:" to "/c"
	rp := regexp.MustCompile("(^[a-zA-Z]):")
	pwds[0] = rp.ReplaceAllString(pwds[0], "/$1")
	pwds[0] = strings.ToLower(pwds[0])

	// Join path using "/"
	return strings.Join(pwds, "/")
}

func parseRepo(str string) (owner, repo string, err error) {
	var parts = strings.Split(str, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid or missing repository name. " +
			"Must use octocat/hello-world format.")
		return
	}
	owner = parts[0]
	repo = parts[1]
	return
}

func resolvePath(dir string) string {

	// attempts to fiture out the project path
	// based on the GOPATH if it exists
	path := os.Getenv("GOPATH")
	path = filepath.Join(path, "src")
	if filepath.HasPrefix(dir, path) {
		return dir[len(path):]
	}

	// attempts to figure out the project path
	// based on common directory conventions
	indexes := gopathExp.FindStringIndex(dir)
	if len(indexes) != 0 {
		index := indexes[len(indexes)-1]
		index = strings.LastIndex(dir, "/src/")
		return dir[index+5:]
	}

	return ""
}

// readInput reads the plaintext secret from a file
// or stdin if inFile is -
func readInput(inFile string) ([]byte, error) {
	if inFile == "-" {
		return ioutil.ReadAll(os.Stdin)
	} else {
		return ioutil.ReadFile(inFile)
	}
}

var gopathExp = regexp.MustCompile("./src/(github.com/[^/]+/[^/]+|bitbucket.org/[^/]+/[^/]+|code.google.com/[^/]+/[^/]+)")

// // getRepoPath checks the source codes absolute path
// // on the host operating system in an attempt
// // to correctly determine the code's package path. This
// // is Go-specific, since Go code must exist in
// // $GOPATH/src/github.com/{owner}/{name}
// func getRepoPath(dir string) (path string, ok bool) {
// 	// let's get the package directory based
// 	// on the path in the host OS
// 	indexes := gopathExp.FindStringIndex(dir)
// 	if len(indexes) == 0 {
// 		return
// 	}

// 	index := indexes[len(indexes)-1]

// 	// if the dir is /home/ubuntu/go/src/github.com/foo/bar
// 	// the index will start at /src/github.com/foo/bar.
// 	// We'll need to strip "/src/" which is where the
// 	// magic number 5 comes from.
// 	index = strings.LastIndex(dir, "/src/")
// 	return dir[index+5:], true
// }

// // getParamMap returns a map of enivronment variables that
// // should be injected into the .drone.yml
// func getParamMap(prefix string) map[string]string {
// 	envs := map[string]string{}

// 	for _, item := range os.Environ() {
// 		env := strings.SplitN(item, "=", 2)
// 		if len(env) != 2 {
// 			continue
// 		}

// 		key := env[0]
// 		val := env[1]
// 		if strings.HasPrefix(key, prefix) {
// 			envs[strings.TrimPrefix(key, prefix)] = val
// 		}
// 	}
// 	return envs
// }

// // prints the time as a human readable string
// func humanizeDuration(d time.Duration) string {
// 	if seconds := int(d.Seconds()); seconds < 1 {
// 		return "Less than a second"
// 	} else if seconds < 60 {
// 		return fmt.Sprintf("%d seconds", seconds)
// 	} else if minutes := int(d.Minutes()); minutes == 1 {
// 		return "About a minute"
// 	} else if minutes < 60 {
// 		return fmt.Sprintf("%d minutes", minutes)
// 	} else if hours := int(d.Hours()); hours == 1 {
// 		return "About an hour"
// 	} else if hours < 48 {
// 		return fmt.Sprintf("%d hours", hours)
// 	} else if hours < 24*7*2 {
// 		return fmt.Sprintf("%d days", hours/24)
// 	} else if hours < 24*30*3 {
// 		return fmt.Sprintf("%d weeks", hours/24/7)
// 	} else if hours < 24*365*2 {
// 		return fmt.Sprintf("%d months", hours/24/30)
// 	}
// 	return fmt.Sprintf("%f years", d.Hours()/24/365)
// }
