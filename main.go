package main

import (
	"github.com/raito-io/cli/cmd"
	"os"
)

var (
	version = "dev"
	date = ""
)

func main() {
	cmd.Execute(buildVersion(version, date), os.Args[1:], os.Exit)
}

func buildVersion(version, date string) string {
	return version + " ("+date+")"
}
