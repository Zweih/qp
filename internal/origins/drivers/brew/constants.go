package brew

const (
	formulaCachePath = "/Homebrew/api/formula.jws.json"
	caskCachePath    = "/Homebrew/api/cask.jws.json"

	receiptName = "INSTALL_RECEIPT.json"

	armMacPrefix       = "/opt/homebrew"
	x86MacPrefix       = "/usr/local"
	x86MacDetectPrefix = "/usr/local/Homebrew"
	linuxPrefix        = "/home/linuxbrew/.linuxbrew"

	cellarSubPath   = "Cellar"
	caskroomSubpath = "Caskroom"
	binSubPath      = "bin"

	typeFormula = "formula"
	typeCask    = "cask"

	dotBrew = ".brew"
	dotRuby = ".rb"

	fieldDesc     = "desc "
	fieldHomepage = "homepage "
	fieldLicense  = "license "

	fieldAllOf = "all_of"
	fieldAnyOf = "any_of"

	spaceAnd = " AND "
	trimAnd  = "AND"
	trimOr   = "OR"

	openBracket  = "["
	closeBracket = "]"
	openParen    = "("
	closeParen   = ")"
	openCurly    = "{"
	closeCurly   = "}"
)
