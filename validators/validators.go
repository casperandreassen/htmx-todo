package validators

import (
	"github.com/go-playground/validator/v10"
	"strconv"
)

var Mod11 validator.Func = func(fl validator.FieldLevel) bool {
	number := fl.Field().String()
	chars := []rune(number)
	sum := 0
	weight := 2
	for i := len(chars) - 2; i >= 0; i-- {
		if num, err := strconv.Atoi(string(chars[i])); err == nil {
			sum += num * weight
			if weight == 7 {
				weight = 2
			} else {
				weight++
			}
		}
	}
	control := 0
	modulo := sum % 11
	if modulo != 0 {
		control = 11 - modulo
	}

	if num, err := strconv.Atoi(string(chars[len(chars)-1])); err == nil {
		if control == num {
			return true
		} else {
			return false
		}
	}
	return false
}
