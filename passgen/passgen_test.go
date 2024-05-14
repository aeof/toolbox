package passgen

import "testing"

func assertLength(t *testing.T, s string, expectedLength int) {
	if len(s) != expectedLength {
		t.Errorf("expect string length %d but %s is of length %d", expectedLength, s, len(s))
	}
}

func TestOneCharset(t *testing.T) {
	t.Run("allDigit", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s := GeneratePassword(i, AllowDigit)
			assertLength(t, s, i)
			for _, ch := range s {
				if ch < '0' || ch > '9' {
					t.Errorf("string %s contains character %c that is not digit\n", s, ch)
				}
			}
		}
	})

	t.Run("allLower", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s := GeneratePassword(i, AllowLower)
			assertLength(t, s, i)
			for _, ch := range s {
				if ch < 'a' || ch > 'z' {
					t.Errorf("string %s contains character %c that is not lowercase letter\n", s, ch)
				}
			}
		}
	})

	t.Run("allUpper", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s := GeneratePassword(i, AllowUpper)
			assertLength(t, s, i)
			for _, ch := range s {
				if ch < 'A' || ch > 'Z' {
					t.Errorf("string %s contains character %c that is not uppercase letter\n", s, ch)
				}
			}
		}
	})

	t.Run("allSymbol", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s := GeneratePassword(i, AllowSymbol)
			assertLength(t, s, i)
			for _, ch := range s {
				if ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch < '!' {
					t.Errorf("string %s contains non-symbol character %c\n", s, ch)
				}
			}
		}
	})
}

func TestMixedCharset(t *testing.T) {
	t.Run("digitAndLower", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			s := GeneratePassword(i, AllowDigit|AllowLower)
			assertLength(t, s, i)
			for _, ch := range s {
				if !(('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'z')) {
					t.Errorf("string %s contains charater %c that is not digit or symbol\n", s, ch)
				}
			}
		}
	})
}
