package file

import (
	"regexp"
	"time"
)

var nonAlphaNumeric = regexp.MustCompile("[^a-zA-Z0-9]+")

func GetFileNameFromName(name string) string {
	return nonAlphaNumeric.ReplaceAllString(name, "")
}

func CreateUniqueFileNameForTarget(target, step, ext string) string {
	name := GetFileNameFromName(target)
	t := time.Now().Format("2006-01-02T15-04-05.000")

	return name + "-" + step + "-" + t + "." + ext
}
