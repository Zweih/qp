package deb

const (
	fieldStatus        = "Status"
	fieldInstalledSize = "Installed-Size"
	fieldPackage       = "Package"
	fieldVersion       = "Version"
	fieldArchitecture  = "Architecture"
	fieldDescription   = "Description"
	fieldHomepage      = "Homepage"
	fieldMaintainer    = "Maintainer"
	fieldConflicts     = "Conflicts"
	fieldBreaks        = "Breaks"
	fieldReplaces      = "Replaces"
	fieldDepends       = "Depends"
	fieldPreDepends    = "Pre-Depends"
	fieldRecommends    = "Recommends"
	fieldSuggests      = "Suggests"
	fieldProvides      = "Provides"
	fieldPriority      = "Priority"
	fieldEssential     = "Essential"
	fieldEnhances      = "Enhances"

	dpkgPath = "/var/lib/dpkg/status"

	installReasonPath = "/var/lib/apt/extended_states"
	packagePrefix     = "Package:"
	autoInstallPrefix = "Auto-Installed:"

	licensePath     = "/usr/share/doc"
	licenseFileName = "copyright"
	filesPrefix     = "Files:"
	licensePrefix   = "License:"
	unknownLicense  = "unknown"

	pkgModRoot = "/var/lib/dpkg/info"
	listExt    = ".list"
	md5SumsExt = ".md5sums"
)
