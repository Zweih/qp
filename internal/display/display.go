package display

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"yaylog/internal/consts"
	"yaylog/internal/pkgdata"

	"golang.org/x/term"
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
		return consts.DefaultTerminalWidth // default width if unable to detect
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

func PrintTable(pkgs []pkgdata.PackageInfo, showFullTimestamp bool, columnNames []string) {
	dateFormat := consts.DateOnlyFormat

	if showFullTimestamp {
		dateFormat = consts.DateTimeFormat
	}

	manager.printTable(pkgs, dateFormat, columnNames)
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
func (o *OutputManager) printTable(
	packages []pkgdata.PackageInfo,
	dateFormat string,
	columnNames []string,
) {
	o.clearProgress()
	columns := []Column{}

	for _, columnName := range columnNames {
		columns = append(columns, GetColumnByName(columnName))
	}

	ctx := displayContext{DateFormat: dateFormat}

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 8, 2, ' ', 0)

	renderHeaders(w, columns)

	for _, pkg := range packages {
		renderRows(w, pkg, columns, ctx)
	}

	w.Flush()
	o.write(buffer.String())
}

func renderHeaders(w *tabwriter.Writer, columns []Column) {
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Header
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))
}

func renderRows(w *tabwriter.Writer, pkg PackageInfo, columns []Column, ctx displayContext) {
	row := make([]string, len(columns))
	for i, col := range columns {
		row[i] = col.Getter(pkg, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}
