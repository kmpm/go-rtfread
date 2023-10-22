package rtfread

func ishex(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// Check if character is an alphabet
func isalpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// Check if character is a digit
func isdigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
