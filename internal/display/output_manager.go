package display

import (
	"fmt"
	"os"
	"strings"
	"sync"
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

func RenderTable(
	pkgs []pkgdata.PkgInfo,
	fields []consts.FieldType,
	showFullTimestamp bool,
	hasNoHeaders bool,
) {
	manager.renderTable(pkgs, fields, showFullTimestamp, hasNoHeaders)
}

func RenderJson(pkgs []pkgdata.PkgInfo, fields []consts.FieldType) {
	manager.renderJson(pkgs, fields)
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
