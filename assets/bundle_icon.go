//go:build ignore

// bundle_icon.go reads icon.png and generates a Go source file with the icon
// embedded as a byte slice, usable as a fyne.StaticResource.
// Run: go run assets/bundle_icon.go
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("assets/icon.png")
	if err != nil {
		panic(err)
	}

	out, err := os.Create("assets/icon_resource.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fmt.Fprintln(out, "package assets")
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import "fyne.io/fyne/v2"`)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "// IconResource is the pixel-art peach icon for Moment.\n")
	fmt.Fprintf(out, "var IconResource = &fyne.StaticResource{\n")
	fmt.Fprintf(out, "\tStaticName:    \"icon.png\",\n")
	fmt.Fprintf(out, "\tStaticContent: []byte{")
	for i, b := range data {
		if i > 0 {
			fmt.Fprint(out, ", ")
		}
		if i%16 == 0 {
			fmt.Fprint(out, "\n\t\t")
		}
		fmt.Fprintf(out, "0x%02x", b)
	}
	fmt.Fprintln(out, ",")
	fmt.Fprintln(out, "\t},")
	fmt.Fprintln(out, "}")
}
