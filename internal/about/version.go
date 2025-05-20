package about

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func PrintVersionInfo() {
	fmt.Printf(
		`  ___   _____
 / __ \/\  __ \
/\ \_\ \ \ \_\ \
\ \___  \ \  __/
 \/___/\ \ \ \/
      \ \_\ \_\
       \/_/\/_/     %s

qp - query packages
https://github.com/Zweih/qp

Version: %s
Commit:  %s
Built:   %s

Copyright (c) 2024-2025 Fernando Nuñez
License: GPLv3-only <https://www.gnu.org/licenses/gpl-3.0.html>
This is free software: you are free to change and redistribute it under GPLv3-only.
There is NO WARRANTY, to the extent permitted by law.

Proprietary redistribution or ML/LLM ingestion requires a separate license.

Author: Fernando Nuñez
`, Version, Version, Commit, Date)
}
