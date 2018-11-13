package validates

import (
	"bytes"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

type Validator struct {
	*validator.Validate
	ut *ut.UniversalTranslator
}

type ValTranslationsErrors struct {
	validator.ValidationErrorsTranslations
}

func NewValidator() *Validator {
	uni := ut.New(en.New(), zh.New())
	transEN, _ := uni.GetTranslator("en")
	transZH, _ := uni.GetTranslator("zh")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, transEN)
	zh_translations.RegisterDefaultTranslations(validate, transZH)
	return &Validator{validate, uni}
}

func IsValidateError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func (v *Validator) GetTranslatedError(err error, locale string) error {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		if locale == "" {
			locale = "en"
		} else if locale == "zh-CN" {
			locale = "zh"
		}
		trans, found := v.ut.GetTranslator(locale)
		if !found {
			trans, _ = v.ut.GetTranslator("en")
		}
		terr := errs.Translate(trans)
		return &ValTranslationsErrors{terr}
	}
	return err
}

func (e *ValTranslationsErrors) Error() string {
	buff := bytes.NewBufferString("")

	for key, val := range e.ValidationErrorsTranslations {
		buff.WriteString(key)
		buff.WriteString(": ")
		buff.WriteString(val)
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
