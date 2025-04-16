package about

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func PrintVersionInfo() {
	fmt.Println()
	fmt.Println("  ___   _____    ")
	fmt.Println(" / __ \\/\\  __ \\ ")
	fmt.Println("/\\ \\_\\ \\ \\ \\_\\ \\")
	fmt.Println("\\ \\___  \\ \\  __/")
	fmt.Println(" \\/___/\\ \\ \\ \\/ ")
	fmt.Println("      \\ \\_\\ \\_\\ ")
	fmt.Println("       \\/_/\\/_/     ", Version)
	fmt.Println()
	fmt.Println("qp - query packages")
	fmt.Println("https://github.com/Zweih/qp")
	fmt.Println()
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Commit:  %s\n", Commit)
	fmt.Printf("Built:   %s\n", Date)
	fmt.Println()
	fmt.Println("Copyright (c) 2024–2025 Fernando Nunez")
	fmt.Println("License GPLv3-only <https://www.gnu.org/licenses/gpl-3.0.html>")
	fmt.Println("This is free software: you are free to change and redistribute it under the GPL.")
	fmt.Println("There is NO WARRANTY, to the extent permitted by law.")
	fmt.Println()
	fmt.Println("Proprietary redistribution or ML/LLM ingestion requires a separate license.")
	fmt.Println()
	fmt.Println("Author: Fernando Nunez")
}
