package helpers

func IsEmpty(value string) bool {
	if len(value) == 0 {
		return true
	}
	return false
}
