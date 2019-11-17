package models

// Contains checks whether or not the array of strings contain the string
func Contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}

	return false
}
