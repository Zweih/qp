package origins

import (
	"qp/api/driver"
	"qp/internal/origins/drivers/brew"
	"qp/internal/origins/drivers/deb"
	"qp/internal/origins/drivers/opkg"
	"qp/internal/origins/drivers/pacman"
	"qp/internal/origins/drivers/pipx"
	"qp/internal/origins/drivers/rpm"
)

var registeredDrivers = []driver.Driver{
	&opkg.OpkgDriver{},
	&deb.DebDriver{},
	&brew.BrewDriver{},
	&pacman.PacmanDriver{},
	&pipx.PipxDriver{},
	&rpm.RpmDriver{},
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
