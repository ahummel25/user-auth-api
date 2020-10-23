package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	enLocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/src/user-auth-api/service"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	en := enLocale.New()
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
	if err = validate.Var(val, constraint); err != nil {
		var transErr error
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			fieldName := *graphql.GetPathContext(ctx).Field
			transErr = fmt.Errorf("%s%+v", fieldName, validationErrors[0].Translate(trans))
		}
		service.SetClientError(service.BAD_REQUEST)
		return val, transErr
	}
	return val, nil
}

// ValidateAddTranslation is a function used for adding a custom validation message
func ValidateAddTranslation(tag string, message string) {
	registerFunc := func(utTrans ut.Translator) error {
		return utTrans.Add(tag, message, true) // see universal-translator for details
	}
	translationFunc := func(utTrans ut.Translator, fe validator.FieldError) string {
		t, _ := utTrans.T(tag, fe.Field())
		return t
	}
	validate.RegisterTranslation(tag, trans, registerFunc, translationFunc)
}
