package consts

const (
	OriginBrew    = "brew"
	OriginDeb     = "deb"
	OriginFlatpak = "flatpak"
	OriginNpm     = "npm"
	OriginOpkg    = "opkg"
	OriginPacman  = "pacman"
	OriginPipx    = "pipx"
	OriginRpm     = "rpm"
	OriginSnap    = "snap"
)

var ValidOrigins = []string{
	OriginBrew,
	OriginDeb,
	OriginFlatpak,
	OriginNpm,
	OriginOpkg,
	OriginPacman,
	OriginPipx,
	OriginRpm,
	OriginSnap,
}
