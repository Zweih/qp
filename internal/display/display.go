package display

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
)

const (
	dateOnlyFormat = "2006-01-02"
	dateTimeFormat = "2006-01-02 15:04:05"
)

const (
	KB = 1024
	MB = KB * KB
	GB = MB * MB
)

type OutputManager struct {
	mu             sync.Mutex
	progressActive bool
	lastMsgLength  int
	terminalWidth  int
}

var manager = newOutputManager()

func newOutputManager() *OutputManager {
	width := getTerminalWidth()

	return &OutputManager{terminalWidth: width}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // default width if unable to detect
	}

	return width
}

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

func PrintTable(pkgs []pkgdata.PackageInfo, showFullTimestamp bool) {
	dateFormat := dateOnlyFormat

	if showFullTimestamp {
		dateFormat = dateTimeFormat
	}

	manager.printTable(pkgs, dateFormat)
}

func (o *OutputManager) write(msg string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	fmt.Print(msg)
}

func (o *OutputManager) printProgress(phase string, progress int, description string) {
	o.progressActive = true

	msg := o.formatProgessMsg(phase, progress, description)
	o.clearPrevMsg(len(msg))

	o.write("\r\033[K" + msg)
	o.lastMsgLength = len(msg)
}

func (o *OutputManager) clearProgress() {
	if o.progressActive {
		o.clearPrevMsg(0)
		o.progressActive = false
	}
}

func (o *OutputManager) formatProgessMsg(phase string, progress int, description string) string {
	msg := fmt.Sprintf("[%s] %d%% - %s", phase, progress, description)

	if len(msg) > o.terminalWidth {
		msg = msg[:o.terminalWidth-1] // truncate message to fit terminal
	}

	return msg
}

func (o *OutputManager) clearPrevMsg(newMsgLength int) {
	if o.lastMsgLength > newMsgLength {
		clearSpace := strings.Repeat(" ", o.lastMsgLength)
		o.write("\r" + clearSpace + "\r")
	}
}

// displays data in tab format
func (o *OutputManager) printTable(pkgs []pkgdata.PackageInfo, dateFormat string) {
	o.clearProgress()

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 8, 2, ' ', 0)

	fmt.Fprintln(w, "DATE\tNAME\tREASON\tSIZE")

	for _, pkg := range pkgs {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			pkg.Timestamp.Format(dateFormat),
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
