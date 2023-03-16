package utils

import (
	"fmt"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

func SetupValidator() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err error
		if err = value.RegisterValidation("maxteams", maxTeams); err != nil {
			log.Fatal(err)
		}
		if err = value.RegisterValidation("checktime", checkTime); err != nil {
			log.Fatal(err)
		}
		if err = value.RegisterValidation("bracket_type", bracketType); err != nil {
			log.Fatal(err)
		}
	}
}

func ValidateErrors(err error) []string {
	var output []string
	if errV, ok := err.(validator.ValidationErrors); ok {
		output = make([]string, len(errV))
		for i, item := range errV {
			output[i] = fmt.Sprintf(
				"failure on field validation: %s",
				item.Field())
		}
	}
	return output
}

var maxTeams validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(int); ok {
		if value%2 == 0 {
			return true
		}
	}
	return false
}

var checkTime validator.Func = func(fl validator.FieldLevel) bool {
	var result bool
	if value, ok := fl.Field().Interface().(time.Time); ok {
		timeNow := time.Now().UTC()
		result = value.UTC().After(timeNow)
	}
	return result
}

var bracketType validator.Func = func(fl validator.FieldLevel) bool {
	var result bool
	if value, ok := fl.Field().Interface().(model.BracketType); ok {
		for _, item := range model.AllBracketType() {
			if value == item {
				result = true
				break
			}
		}
	}
	return result
}
