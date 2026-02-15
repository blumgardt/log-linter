package rules

func ViolatesEnglishOnlyASCII(msg string) bool {
	for _, r := range msg {
		if r > 127 {
			return true
		}
	}
	return false
}
