package version

var v = ""

func SetVersion(version, date string) {
	v = version + " (" + date + ")"
}

func GetVersion() string {
	return v
}
