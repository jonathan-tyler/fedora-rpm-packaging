package commandline

import "strings"

func Quote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `"'"'`) + "'"
}
