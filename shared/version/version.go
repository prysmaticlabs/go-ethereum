// Package version executes and returns the version string
// for the currently running process.
package version

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// The value of these vars are set through linker options.
var gitCommit = "Local build"
var buildDate = "Moments ago"
var gitTag = "Unknown"

// GetVersion returns the version string of this build.
func GetVersion() string {
	// if doing a local build, these values are not interpolated
	if gitCommit == "{STABLE_GIT_COMMIT}" {
		commit, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			log.Println(err)
		} else {
			gitCommit = strings.TrimRight(string(commit), "\r\n")
		}
	}
	if buildDate == "{DATE}" {
		now := time.Now().Format(time.RFC3339)
		buildDate = now
	}
	return fmt.Sprintf("Prysm/%s/%s. Built at: %s", gitTag, gitCommit, buildDate)
}
