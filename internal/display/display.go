package display

import (
  "fmt"
  "os"
  "text/tabwriter"
  "yaylog/internal/pkgdata"
)

// displays data in tab format
func PrintTable(pkgs []pkgdata.PackageInfo) {
  w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
  fmt.Fprintln(w, "DATE\tNAME\tREASON")

  for _, pkg := range pkgs {
    fmt.Fprintf(
      w,
      "%s\t%s\t%s\n",
      pkg.Timestamp.Format("2006-01-02 15:04:05"),
      pkg.Name,
      pkg.Reason,
    )
  }

  w.Flush()
}
