package pipx

const (
	defaultVenvRoot = ".local/pipx/venvs"
	otherVenvRoot   = ".local/share/pipx/venvs"

	fieldName        = "Name"
	fieldVersion     = "Version"
	fieldSummary     = "Summary"
	fieldLicense     = "License"
	fieldHomepage    = "Home-page"
	fieldProjectUrl  = "Project-Url"
	subfieldHomepage = "Homepage"
	fieldTag         = "Tag"

	anyArch       = "any"
	universalArch = "universal"

	dotDistInfo = ".dist-info"

	pipxHomeEnv = "PIPX_HOME"
	homeEnv     = "HOME"
)
