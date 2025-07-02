package origins

import (
	"qp/api/driver"
	"qp/internal/origins/drivers/brew"
	"qp/internal/origins/drivers/deb"
	"qp/internal/origins/drivers/flatpak"
	"qp/internal/origins/drivers/npm"
	"qp/internal/origins/drivers/opkg"
	"qp/internal/origins/drivers/pacman"
	"qp/internal/origins/drivers/pipx"
	"qp/internal/origins/drivers/rpm"
	"qp/internal/origins/drivers/snap"
)

var registeredDrivers = []driver.Driver{
	&brew.BrewDriver{},
	&deb.DebDriver{},
	&flatpak.FlatpakDriver{},
	&opkg.OpkgDriver{},
	&npm.NpmDriver{},
	&pacman.PacmanDriver{},
	&pipx.PipxDriver{},
	&rpm.RpmDriver{},
	&snap.SnapDriver{},
}

func AvailableDrivers() []driver.Driver {
	var detected []driver.Driver
	for _, driver := range registeredDrivers {
		if driver.Detect() {
			detected = append(detected, driver)
		}
	}

	return detected
}
