package opkg

const (
	opkgStatusPath = "/usr/lib/opkg/status"

	fieldPackage       = "Package"
	fieldVersion       = "Version"
	fieldArchitecture  = "Architecture"
	fieldDepends       = "Depends"
	fieldProvides      = "Provides"
	fieldConflicts     = "Conflicts"
	fieldInstalledTime = "Installed-Time"
	fieldStatus        = "Status"
	fieldAutoInstalled = "Auto-Installed"
	fieldEssential     = "Essential"

	opkgInfoRoot = "/usr/lib/opkg/info"

	fieldLicense         = "License"
	fieldDescription     = "Description"
	fieldInstalledSize   = "Installed-Size"
	fieldSourceDateEpoch = "SourceDateEpoch"
)
