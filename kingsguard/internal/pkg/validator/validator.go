package validator

import (
	"fmt"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func Password() validation.Rule {
	// TODO: написать более оптимальную проверку через регулярку
	return validation.By(func(value interface{}) error {
		str, _ := value.(string)

		if len(str) < 8 || len(str) > 20 {
			return fmt.Errorf("password length must be greater than 8 and less than 20")
		}

		hasLetter := false
		hasNonLetter := false

		for _, r := range str {
			if unicode.IsLetter(r) {
				hasLetter = true
			} else {
				hasNonLetter = true
			}
		}

		if !hasLetter {
			return fmt.Errorf("password must have at least one letter")
		}
		if !hasNonLetter {
			return fmt.Errorf("password must have at least one symbol (not the letter)")
		}

		return nil
	})
}
