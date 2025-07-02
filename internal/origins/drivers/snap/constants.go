package snap

const (
	snapRoot = "/snap"

	networkUnix         = "unix"
	snapdSocket         = "/run/snapd.socket"
	snapLocalHost       = "http://localhost/v2/snaps"
	connectionLocalHost = "http://localhost/v2/connections"

	typeOs    = "os"
	typeSnapd = "snapd"

	binDir   = "/var/lib/snapd/snap/bin"
	snapsDir = "/var/lib/snapd/snaps"
	dotSnap  = ".snap"

	interfaceContent = "content"

	timestampFormat = "2006-01-02T15:04:05.999999999-07:00"
)
