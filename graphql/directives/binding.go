package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, trans)
}

// Binding implements the binding directive function and handles any field or input validation errors
func Binding(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}
	fieldName := *graphql.GetPathContext(ctx).Field
	err = validate.Var(val, constraint)
	if err != nil {
		var transErr error
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			transErr = fmt.Errorf("%s%+v", fieldName, validationErrors[0].Translate(trans))
		}
		return val, transErr
	}

	return val, nil
}

// ValidateAddTranslation is a function used for adding a custom validation message
func ValidateAddTranslation(tag string, message string) {
	validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, message, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}
