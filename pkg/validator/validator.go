package validator

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	//Register custom validations
	_ = v.RegisterValidation("price", validatePrice)
	_ = v.RegisterValidation("quantity", validateQuantity)

	return &Validator{
		validate: v,
	}
}

func (v *Validator) Validate(ctx context.Context, i interface{}) error {
	return v.validate.StructCtx(ctx, i)
}

func (v *Validator) ValidateVar(ctx context.Context, field interface{}, tag string) error {
	return v.validate.VarCtx(ctx, field, tag)
}

func (v *Validator) ExtractValidationErrors(err error) []ValidationError {
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Message: v.generateValidationMessage(err),
		})
	}
	return validationErrors
}

func (v *Validator) generateValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Failed ! This field is required, please fill it"
	case "email":
		return "Failed ! Invalid email"
	case "min":
		return fmt.Sprintf("Failed ! Value should be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("Failed ! value should be at most %s", err.Param())
	case "price":
		return "Failed ! Price must be greater than 0"
	case "quantity":
		return "Failed ! Quantity must be 0 or greater"
	default:
		return fmt.Sprintf("Failed on %s validation", err.Tag())
	}
}

func validatePrice(fl validator.FieldLevel) bool {
	price, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}
	return price > 0
}

func validateQuantity(fl validator.FieldLevel) bool {
	quantity, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	return quantity >= 0
}

func (v *Validator) ValidateID(ctx context.Context, id string) error {
	return v.ValidateVar(ctx, id, "required,uuid")
}

func (v *Validator) ValidatePagination(ctx context.Context, page, pageSize int) error {
	if err := v.ValidateVar(ctx, page, "min=1"); err != nil {
		return fmt.Errorf("invalid page number: %w", err)
	}
	if err := v.ValidateVar(ctx, pageSize, "min=1,max=100"); err != nil {
		return fmt.Errorf("invalid page size: %w", err)
	}
	return nil
}
