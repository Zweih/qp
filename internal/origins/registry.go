package origins

import (
	"qp/interfaces"
	"qp/internal/origins/apt"
	"qp/internal/origins/pacman"
)

var registeredDrivers = []interfaces.Driver{
	&apt.AptDriver{},
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
