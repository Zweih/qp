package shared

import (
	"os"
	"strings"
	"sync"
)

const (
	dockerEnvPath = "/.dockerenv"
	cGroupPath    = "/proc/1/cgroup"
	rootDir       = "/"
)

var (
	creationTimeReliable *bool
	once                 sync.Once
)

// for one-off requests
func GetCreationTime(path string) (int64, bool, error) {
	once.Do(func() {
		reliable := checkCreationTimeReliability()
		creationTimeReliable = &reliable
	})

	if !*creationTimeReliable {
		return 0, false, nil
	}

	return getBirthTime(path)
}

// for use to preemptively avoid looping
func IsCreationTimeReliable() bool {
	once.Do(func() {
		reliable := checkCreationTimeReliability()
		creationTimeReliable = &reliable
	})

	return *creationTimeReliable
}

func checkCreationTimeReliability() bool {
	if detectDocker() {
		return false
	}

	if _, reliable, err := getBirthTime(rootDir); err != nil || !reliable {
		return false
	}

	return true
}

func detectDocker() bool {
	if _, err := os.Stat(dockerEnvPath); err == nil {
		return true
	}

	if data, err := os.ReadFile(cGroupPath); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") || strings.Contains(content, "/docker/") {
			return true
		}
	}

	return false
}
