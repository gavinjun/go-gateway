package string_util

import "strings"

func StrCompareIgnoreLowerOrUpper(a,b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}
