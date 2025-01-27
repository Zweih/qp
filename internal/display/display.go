package display

import (
	"fmt"
	"os"
	"text/tabwriter"
	"yaylog/internal/pkgdata"
)

const (
	KB = 1024
	MB = KB * KB
	GB = MB * MB
)

// displays data in tab format
func PrintTable(pkgs []pkgdata.PackageInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tNAME\tREASON\tSIZE")

	for _, pkg := range pkgs {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			pkg.Timestamp.Format("2006-01-02 15:04:05"),
			pkg.Name,
			pkg.Reason,
			formatSize(pkg.Size),
		)
	}

	w.Flush()
}

func formatSize(size int64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
