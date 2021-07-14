package util

// element が source にあるか
func Contains(source []string, element string) bool {
	for _, item := range source {
		if item == element {
			return true
		}
	}
	return false
}
