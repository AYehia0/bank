// custom validators
package api

import (
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// check currency support
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
