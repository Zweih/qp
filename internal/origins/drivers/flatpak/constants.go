package flatpak

const (
	systemInstallDir = "/var/lib/flatpak/"
	userInstallDir   = ".local/share/flatpak/"

	typeApp        = "app"
	typeRuntime    = "runtime"
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
	sectionExtensionOf = "ExtensionOf"
	fieldName          = "name"
	fieldVersion       = "version"
	fieldRef           = "ref"
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
