package display

import (
	"bytes"
	"encoding/json"
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
	manager.writeLine(msg)
}

func PrintProgress(phase string, progress int, description string) {
	manager.printProgress(phase, progress, description)
}

func ClearProgress() {
	manager.clearProgress()
}

func PrintTable(pkgs []pkgdata.PackageInfo, columnNames []string, showFullTimestamp bool, hasNoHeaders bool) {
	dateFormat := consts.DateOnlyFormat

	if showFullTimestamp {
		dateFormat = consts.DateTimeFormat
	}

	manager.printTable(pkgs, dateFormat, columnNames, hasNoHeaders)
}

func PrintJson(pkgs []pkgdata.PackageInfo, columnNames []string) {
	manager.printJson(pkgs, columnNames)
}

func (o *OutputManager) write(msg string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	fmt.Print(msg)
}

func (o *OutputManager) writeLine(msg string) {
	o.write(msg + "\n")
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
	hasNoHeaders bool,
) {
	o.clearProgress()
	ctx := displayContext{DateFormat: dateFormat}

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 8, 2, ' ', 0)

	if !hasNoHeaders {
		renderHeaders(w, columnNames)
	}

	for _, pkg := range packages {
		renderRows(w, pkg, columnNames, ctx)
	}

	w.Flush()
	o.write(buffer.String())
}

func renderHeaders(w *tabwriter.Writer, columnNames []string) {
	headers := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		headers[i] = columnHeaders[columnName]
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))
}

func renderRows(w *tabwriter.Writer, pkg pkgdata.PackageInfo, columnNames []string, ctx displayContext) {
	row := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		row[i] = GetColumnTableValue(pkg, columnName, ctx)
	}

	fmt.Fprintln(w, strings.Join(row, "\t"))
}

func (o *OutputManager) printJson(pkgs []pkgdata.PackageInfo, columnNames []string) {
	filteredPackages := make([]pkgdata.PackageInfoJson, len(pkgs))
	for i, pkg := range pkgs {
		filteredPackages[i] = GetColumnJsonValues(pkg, columnNames)
	}

	jsonOutput, err := json.MarshalIndent(filteredPackages, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}
