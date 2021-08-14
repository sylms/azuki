package util

import (
	"strings"
)

// element が source にあるか
func Contains(source []string, element string) bool {
	for _, item := range source {
		if item == element {
			return true
		}
	}
	return false
}

// 与えられた text を半角・全角スペースで区切る
func SplitSpace(text string) []string {
	// f := func(c rune) bool {
	// 	return unicode.IsSpace(c)
	// }

	// splitted := strings.FieldsFunc(text, f)
	// return splitted

	return strings.Fields(text)
}
