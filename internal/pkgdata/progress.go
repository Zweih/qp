package pkgdata

type ProgressMessage struct {
	Phase       string
	Progress    int
	Description string
}

type ProgressReporter func(current int, total int, phase string)
