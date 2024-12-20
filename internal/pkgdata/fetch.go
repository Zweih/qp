package pkgdata

import (
	"time"
)

func FetchPackages() ([]PackageInfo, error) {
  packages := []PackageInfo{
    {
      Timestamp: time.Now().Add(-time.Hour * 24 * 3), // 3 days ago
      Name:      "foo",
      Reason:    "explicit",
    },
    {
      Timestamp: time.Now().Add(-time.Hour * 24), // 1 day ago
      Name:      "bar",
      Reason:    "dependency",
    },
    {
      Timestamp: time.Now(), // now
      Name:      "baz",
      Reason:    "explicit",
    },    
  }

  return packages, nil
}
