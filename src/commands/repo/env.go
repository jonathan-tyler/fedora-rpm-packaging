package repo

import "os"

func repoEnvironment() (string, string, string) {
	releasever := getenvDefault("FEDORA_PACKAGE_RELEASEVER", "42")
	basearch := getenvDefault("FEDORA_PACKAGE_BASEARCH", "x86_64")
	gpgKey := os.Getenv("FEDORA_PACKAGE_GPG_KEY")

	return releasever, basearch, gpgKey
}

func getenvDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}