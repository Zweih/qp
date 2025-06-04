package rpm

const (
	defaultRpmRoot = "/var/lib/rpm"
	modernRpmRoot  = "/usr/lib/sysimage/rpm"
	rebuildRpmRoot = "/var/lib/rpmrebuilddb"

	dnfHistoryPath    = "/var/lib/dnf/history.sqlite"
	yumHistoryPath    = "/var/lib/yum/history.sqlite"
	yumHistoryPattern = "/var/lib/yum/history/history-*.sqlite"

	sqliteDbFile   = "rpmdb.sqlite"
	ndbDbFile      = "Packages.db"
	berkeleyDbFile = "Packages"

	groupUnspecified = "Unspecified"

	versionTilde = "~"
	digitZero    = "0"

	dnfReasonDependency = 1
	dnfReasonUser       = 2
	dnfReasonWeakDep    = 4

	dnfActionInstall  = 1 // new package installation
	dnfActionUpgrade  = 6 // installing new version of existing package
	dnfActionUpgraded = 7 // removing old version during upgrade
	dnfActionRemove   = 8

	yumReasonUser = "user"
	yumReasonDep  = "dep"

	yumStateInstall     = "Install"
	yumStateTrueInstall = "True-Install"
	yumStateDepInstall  = "Dep-Install"

	opAnd    = " and "
	opOr     = " or "
	opIf     = " if "
	opUnless = " unless "

	opOpenParen  = '('
	opCloseParen = ')'

	valueOpenParen  = "("
	valueCloseParen = ")"
	valueAnd        = "and"
	valueOr         = "or"
	valueIf         = "if"
	valueUnless     = "unless"

	opGreaterEqual = ">="
	opLessEqual    = "<="
	opGreater      = ">"
	opLess         = "<"
	opEqual        = "="
)

type TokenType int

const (
	tokenPackage TokenType = iota
	tokenAnd
	tokenOr
	tokenIf
	tokenUnless
	tokenWith
	tokenWithout
	tokenLParen
	tokenRParen
	tokenEOF
	tokenError
)
