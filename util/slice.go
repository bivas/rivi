package util

func InStringSlice(slice []string, lookup string) bool {
	for _, item := range slice {
		if item == lookup {
			return true
		}
	}
	return false
}
