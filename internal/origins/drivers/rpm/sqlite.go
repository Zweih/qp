//go:build !nosqlite

package rpm

import (
	_ "github.com/glebarez/sqlite"
)
