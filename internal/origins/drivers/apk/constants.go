package apk

const (
	apkDbPath    = "/lib/apk/db/installed"
	apkWorldPath = "/etc/apk/world"

	fieldChecksum      = "C"
	fieldPackage       = "P"
	fieldVersion       = "V"
	fieldArch          = "A"
	fieldSize          = "S"
	fieldInstalledSize = "I"
	fieldDescription   = "T"
	fieldUrl           = "U"
	fieldLicense       = "L"
	fieldOrigin        = "o"
	fieldMaintainer    = "m"
	fieldBuildTime     = "t"
	fieldCommit        = "c"
	fieldDepends       = "D"
	fieldProvides      = "p"
	fieldInstallIf     = "i"
	fieldReplaces      = "r"
	fieldRepoTag       = "s"
	fieldBroken        = "f"
	fieldDirectory     = "F"
	fieldFile          = "R"
)
