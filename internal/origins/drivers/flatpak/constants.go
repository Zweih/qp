package flatpak

const (
	systemInstallDir = "/var/lib/flatpak"
	userInstallDir   = ".local/share/flatpak"

	appsSubdir     = "app"
	runtimeSubdir  = "runtime"
	repoSubdir     = "repo"
	metadataSubdir = "metadata"

	metadataFile = "metadata"
	activeFile   = "active"

	sectionApplication = "Application"
	sectionRuntime     = "Runtime"
	fieldName          = "name"
	fieldVersion       = "version"
	fieldRuntime       = "runtime"
	fieldBranch        = "branch"
	fieldArch          = "arch"
	fieldSdkExtends    = "sdk-extensions"
	fieldFinishArgs    = "finish-args"
	fieldCommand       = "command"

	fieldExtensionOf = "Extension"
	fieldDirectory   = "directory"
	fieldLocale      = "locale"

	permissionPrefix     = "--"
	sharePrefix          = "--share="
	socketPrefix         = "--socket="
	filesystemPrefix     = "--filesystem="
	devicePrefix         = "--device="
	allowPrefix          = "--allow="
	disallowPrefix       = "--disallow="
	envPrefix            = "--env="
	ownNamePrefix        = "--own-name="
	talkNamePrefix       = "--talk-name="
	systemTalkNamePrefix = "--system-talk-name="

	defaultArch = "x86_64"
)
