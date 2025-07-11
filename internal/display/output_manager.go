package display

import (
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
	"sync"

	"golang.org/x/term"
)

const DefaultTerminalWidth = 80

type OutputManager struct {
	mu            sync.Mutex
	lastMsgLength int
	terminalWidth int
}

var manager = newOutputManager()

func newOutputManager() *OutputManager {
	width := getTerminalWidth()

	return &OutputManager{terminalWidth: width}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return DefaultTerminalWidth // default width if unable to detect
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

func RenderTable(
	pkgPtrs []*pkgdata.PkgInfo,
	fields []consts.FieldType,
	showFullTimestamp bool,
	hasNoHeaders bool,
) {
	manager.renderTable(pkgPtrs, fields, showFullTimestamp, hasNoHeaders)
}

func RenderJSON(pkgPtrs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	manager.renderJSON(pkgPtrs, fields)
}

func RenderKeyValue(pkgs []*pkgdata.PkgInfo, fields []consts.FieldType) {
	manager.renderKeyValue(pkgs, fields)
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
	msg := o.formatProgessMsg(phase, progress, description)
	o.clearPrevMsg(len(msg))

	o.write("\r\033[K" + msg)
	o.lastMsgLength = len(msg)
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
