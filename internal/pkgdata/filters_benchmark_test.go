package pkgdata

import (
	"testing"
)

// Fetch real package data
var realPackages, _ = FetchPackages()

// Benchmark passing entire struct (current approach)
func BenchmarkFilterByNamesWithStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, pkg := range realPackages {
			_ = FilterByNames(pkg, []string{"fire", "fox", "chromium"})
		}
	}
}

// Benchmark passing only string (optimized approach)
func BenchmarkFilterByNamesWithString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, pkg := range realPackages {
			_ = FilterByNamesOptimized(pkg.Name, []string{"fire", "fox", "chromium"})
		}
	}
}

func BenchmarkFilterByNamesWithStringAndIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, pkg := range realPackages {
			_ = FilterByNamesOptimized2(pkg.Name, []string{"fire", "fox", "chromium"})
		}
	}
}
