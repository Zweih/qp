package shared

import (
	"os"
	"strings"
	"sync"
)

const (
	dockerEnvPath = "/.dockerenv"
	cGroupPath    = "/proc/1/cgroup"
)

var (
	inDocker bool
	once     sync.Once
)

// TODO: perhaps we should have this all part of "reliable". we should include checking if inside of docker as a part of reliability, but we only check reliability once.
func InDocker() bool {
	once.Do(func() {
		inDocker = detectDocker()
	})

	return inDocker
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
