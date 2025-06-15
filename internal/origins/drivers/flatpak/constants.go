package flatpak

const (
	systemInstallDir = "/var/lib/flatpak/"
	userInstallDir   = ".local/share/flatpak/"

	appsSubdir     = "app"
	runtimeSubdir  = "runtime"
	repoSubdir     = "repo"
	metadataSubdir = "metadata"

	scopeSystem = "system"
	scopeUser   = "user"

	remotesDir    = "repo/refs/remotes/"
	remoteUnknown = "unknown"

	metadataFile = "metadata"
	activeFile   = "active"

	sectionApplication = "Application"
	sectionRuntime     = "Runtime"
	sectionExtension   = "Extension "
	fieldName          = "name"
	fieldVersion       = "version"
	fieldRuntime       = "runtime"
	fieldBranch        = "branch"
	fieldArch          = "arch"

	fieldDirectory = "directory"
	fieldLocale    = "locale"

	metainfoDir    = "files/share/metainfo/"
	appdataDir     = "files/share/appdata/"
	dotMetainfoXml = ".metainfo.xml"
	dotAppdataXml  = ".appdata.xml"

	applicationsDir = "files/share/applications/"
	dotDesktop      = ".desktop"

	appdataVersion = "appdata-version"
)
