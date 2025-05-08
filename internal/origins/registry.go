package origins

import (
	"qp/interfaces"
	"qp/internal/origins/deb"
	"qp/internal/origins/opkg"
	"qp/internal/origins/pacman"
)

var registeredDrivers = []interfaces.Driver{
	&opkg.OpkgDriver{},
	&deb.DebDriver{},
	&pacman.PacmanDriver{},
}

func AvailableDrivers() []interfaces.Driver {
	var detected []interfaces.Driver
	for _, driver := range registeredDrivers {
		if driver.Detect() {
			detected = append(detected, driver)
		}
	}

	return detected
}
