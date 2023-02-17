package main

import (
	"os"

	"github.com/raito-io/cli/cmd"
	v "github.com/raito-io/cli/internal/version"
)

var (
	version = v.DevVersion.String()
	date    = ""
)

func main() {
	v.SetVersion(version, date)
	cmd.Execute(v.GetVersionString(), os.Args[1:], os.Exit)
}
