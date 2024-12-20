package pkgdata

import "time"

// PackageInfo represents information about a single installed package
type PackageInfo struct {
  Timestamp time.Time
  Name      string
  Reason    string // "explicit" or "dependency"
}

