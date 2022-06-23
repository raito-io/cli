package main

import (
	"os"

	"github.com/raito-io/cli/cmd"
	v "github.com/raito-io/cli/internal/version"
)

var (
	version = "dev"
	date    = ""
)

func main() {
	v.SetVersion(version, date)
	cmd.Execute(v.GetVersion(), os.Args[1:], os.Exit)
}
