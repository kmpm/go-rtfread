package internal

func IsHex(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// Check if character is an alphabet
func IsAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// Check if character is a digit
func IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
