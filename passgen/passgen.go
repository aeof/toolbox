package passgen

import (
	"math/rand"
	"time"
)

/*
"passgen" is a tool to generate strong password of length N, for each character:
1. randomly select a type of character set(digit, uppercase alphabet, lowercase alphabet, symbols)
2. randomly select a character in the character set
*/

// flags to specify the selected charsets for password generation
const (
	AllowDigit = 1 << iota
	AllowLower
	AllowUpper
	AllowSymbol
)

var (
	random *rand.Rand

	symbols = []byte{
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*',
		'+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
		'?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~',
	}
)

func init() {
	// initialize random generator
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

var charsetGenerators = []func() byte{
	generateDigit, generateLower, generateUpper, generateSymbol,
}

func generateDigit() byte {
	return '0' + byte(random.Intn(10))
}

func generateLower() byte {
	return 'a' + byte(random.Intn(26))
}

func generateUpper() byte {
	return 'A' + byte(random.Intn(26))
}

func generateSymbol() byte {
	return symbols[random.Intn(len(symbols))]
}

// GeneratePassword generates password with the specified charsetFlags
// Example: to generate passwords with lowercase letters and digit, just pass `AllowLower|AllowDigit`
func GeneratePassword(length int, charsetFlags int) string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		indexCharsetGenerator := -1
		for indexCharsetGenerator == -1 || charsetFlags&(1<<indexCharsetGenerator) == 0 {
			indexCharsetGenerator = random.Intn(len(charsetGenerators))
		}
		buf[i] = charsetGenerators[indexCharsetGenerator]()
	}
	return string(buf)
}
