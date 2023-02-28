package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var (
	LenPassword      int
	IncludeNumbers   bool
	IncludeLowerCase bool
	IncludeUpperCase bool
	IncludeSymbols   bool
)

var (
	symbols = []byte{
		'!', '@', '#', '$', '%', '^', '&', '*', '(', ')',
		'-', '_', '+', '=', '[', ']', '{', '}', ':', ';',
		'"', '\'', '<', '>', ',', '.', '?', '/', '|', '\\',
	}
	generators = map[string]func() byte{
		// numbers
		"num": func() byte {
			return byte('0' + rand.Intn(10))
		},
		// lowercase characters
		"lower": func() byte {
			return byte('a' + rand.Intn(26))
		},
		// uppercase characters
		"upper": func() byte {
			return byte('A' + rand.Intn(26))
		},
		// symbols
		"symbol": func() byte {
			return symbols[rand.Intn(len(symbols))]
		},
	}
)

var PassGenCmd = &cobra.Command{
	Use:   "passgen",
	Short: "PassGen is a tool for generating website password locally",
	Run: func(cmd *cobra.Command, args []string) {
		generatePassword()
	},
}

func init() {
	{
		PassGenCmd.Flags().BoolVar(&IncludeNumbers, "num", true, "include numbers")
		PassGenCmd.Flags().BoolVar(&IncludeLowerCase, "lower", true, "include lowercase characters")
		PassGenCmd.Flags().BoolVar(&IncludeUpperCase, "upper", true, "include uppercase characters")
		PassGenCmd.Flags().BoolVar(&IncludeSymbols, "symbol", false, "include symbols")
	}
	{
		PassGenCmd.Flags().IntVar(&LenPassword, "len", 16, "password length")
	}

	// TODO: fix deprecated seeding
	rand.Seed(time.Now().Unix())
}

func generatePassword() {
	mapping := map[string]bool{
		"num":    IncludeNumbers,
		"lower":  IncludeLowerCase,
		"upper":  IncludeUpperCase,
		"symbol": IncludeSymbols,
	}

	var generatorSet []func() byte
	for chType, included := range mapping {
		if included {
			generatorSet = append(generatorSet, generators[chType])
		}
	}

	bytes := make([]byte, LenPassword)
	for i := 0; i < LenPassword; i++ {
		bytes[i] = generatorSet[rand.Intn(len(generatorSet))]()
	}
	fmt.Println(string(bytes))
}
