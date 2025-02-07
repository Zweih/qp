package display

import (
	"bytes"
	"fmt"
	"sync"
	"text/tabwriter"
	"yaylog/internal/pkgdata"
)

const (
	KB = 1024
	MB = KB * KB
	GB = MB * MB
)

type OutputManager struct {
	mu             sync.Mutex
	progressActive bool
}

var manager = &OutputManager{}

func Write(msg string) {
	manager.write(msg)
}

func WriteLine(msg string) {
	manager.write(msg + "\n")
}

func PrintProgress(phase string, progress int, description string) {
	manager.printProgress(phase, progress, description)
}

func ClearProgress() {
	manager.clearProgress()
}

func PrintTable(pkgs []pkgdata.PackageInfo) {
	manager.printTable(pkgs)
}

func (o *OutputManager) write(msg string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	fmt.Print(msg)
}

func (o *OutputManager) printProgress(phase string, progress int, description string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.progressActive = true

	fmt.Print("\r\033[K")
	fmt.Printf("\r[%s] %d%% - %s", phase, progress, description)
}

func (o *OutputManager) clearProgress() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.progressActive {
		fmt.Print("\r\033[K")
		o.progressActive = false
	}
}

// displays data in tab format
func (o *OutputManager) printTable(pkgs []pkgdata.PackageInfo) {
	o.clearProgress()

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 8, 2, ' ', 0)

	fmt.Fprintln(w, "DATE\tNAME\tREASON\tSIZE")

	for _, pkg := range pkgs {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			pkg.Timestamp.Format("2006-01-02 15:04:05"),
			pkg.Name,
			pkg.Reason,
			formatSize(pkg.Size),
		)
	}

	w.Flush()

	o.write(buffer.String())
}

func formatSize(size int64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
