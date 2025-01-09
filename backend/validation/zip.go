package validation

func IsValidZipCode(zipCode string) bool {
	if len(zipCode) != 5 {
		return false
	}
	for _, char := range zipCode {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}