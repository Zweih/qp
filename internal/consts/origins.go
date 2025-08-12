package consts

const (
	OriginApk     = "apk"
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
	OriginApk,
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
