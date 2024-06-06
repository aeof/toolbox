package cmd

import (
	"errors"
	"fmt"
	"github.com/aeof/toolbox/passgen"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	defaultPasswordLength = 18
)

var (
	passgenCmd = &cobra.Command{
		Use:   "passgen",
		Short: "Passgen is a simple tool to generate a strong password",
		Long: `Passgen is a simple tool to generate strong password. 
You can specify password length, which defaults to 18.
You can also specify the character sets, which defaults to digits, lowercase and uppercase letters`,
		RunE: RunPassgen,
	}

	// passgen command flags
	allowDigit  bool
	allowLower  bool
	allowUpper  bool
	allowSymbol bool
)

func init() {
	passgenCmd.Flags().BoolVar(&allowDigit, "digit", true, "allow digit characters")
	passgenCmd.Flags().BoolVar(&allowLower, "lower", true, "allow lowercase letters")
	passgenCmd.Flags().BoolVar(&allowUpper, "upper", true, "allow uppercase letters")
	passgenCmd.Flags().BoolVar(&allowSymbol, "symbol", false, "allow symbol letters")
	rootCmd.AddCommand(passgenCmd)
}

func RunPassgen(cmd *cobra.Command, args []string) error {
	lenPassword := defaultPasswordLength
	if len(args) > 0 {
		lenPassword, err := strconv.Atoi(args[0])
		if err != nil || lenPassword <= 0 {
			return errors.New("password length should be positive")
		}
	}

	flagLists := []struct {
		allowed     bool
		passgenFlag int
	}{
		{allowDigit, passgen.AllowDigit},
		{allowLower, passgen.AllowLower},
		{allowUpper, passgen.AllowUpper},
		{allowSymbol, passgen.AllowSymbol},
	}

	var flags int
	for _, flagItem := range flagLists {
		if flagItem.allowed {
			flags |= flagItem.passgenFlag
		}
	}
	fmt.Println(passgen.GeneratePassword(lenPassword, flags))
	return nil
}
