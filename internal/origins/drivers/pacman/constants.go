package pacman

const (
	fieldInstallDate = "%INSTALLDATE%"
	fieldBuildDate   = "%BUILDDATE%"
	fieldName        = "%NAME%"
	fieldSize        = "%SIZE%"
	fieldReason      = "%REASON%"
	fieldVersion     = "%VERSION%"
	fieldArch        = "%ARCH%"
	fieldLicense     = "%LICENSE%"
	fieldPkgBase     = "%BASE%"
	fieldDescription = "%DESC%"
	fieldUrl         = "%URL%"
	fieldValidation  = "%VALIDATION%"
	fieldPackager    = "%PACKAGER%"
	fieldGroups      = "%GROUPS%"
	fieldDepends     = "%DEPENDS%"
	fieldOptDepends  = "%OPTDEPENDS%"
	fieldProvides    = "%PROVIDES%"
	fieldConflicts   = "%CONFLICTS%"
	fieldReplaces    = "%REPLACES%"
	fieldXData       = "%XDATA%"

	subfieldPkgType = "pkgtype"

	pacmanDbDir   = "/var/lib/pacman/local"
	pacmanLogPath = "/var/log/pacman.log"

	etcMachineId = "/etc/machine-id"
	etcHostname  = "/etc/machine-id"
	bootLinux    = "/boot/vmlinuz-linux"
	binPacman    = "/usr/bin/pacman"
	etcPasswd    = "/etc/passwd"
	etcGroup     = "/etc/group"
)
