package utils

import (
	customerModel "gin-dbo/view/customer"
	loginModel "gin-dbo/view/login"
	orderModel "gin-dbo/view/order"

	"github.com/go-playground/validator/v10"
)

var (
	// login
	loginRule = map[string]string{
		"Username": "required",
		"Password": "required",
	}
	createLoginRule = map[string]string{
		"Username": "required",
		"Password": "required",
		"Role":     "required",
	}
	updateLoginRule = map[string]string{
		"Username": "required",
		"Password": "required",
		"Role":     "required",
	}

	// customer
	createCustomerRule = map[string]string{
		"Name": "required",
	}
	updateCustomerRule = map[string]string{
		"Id":   "required",
		"Name": "required",
	}
	deleteCustomerRule = map[string]string{
		"Id": "required",
	}

	// Order
	createOrderRule = map[string]string{
		"Name":       "required",
		"Qty":        "required,min=0",
		"CustomerId": "required",
	}
	updateOrderRule = map[string]string{
		"Id":         "required",
		"Name":       "required",
		"Qty":        "required,min=0",
		"CustomerId": "required",
	}
	deleteOrderRule = map[string]string{
		"Id": "required",
	}
)

func NewValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterStructValidationMapRules(loginRule, loginModel.LoginRequest{})
	validate.RegisterStructValidationMapRules(createLoginRule, loginModel.CreateRequest{})
	validate.RegisterStructValidationMapRules(updateLoginRule, loginModel.UpdateRequest{})
	validate.RegisterStructValidationMapRules(createCustomerRule, customerModel.CreateRequest{})
	validate.RegisterStructValidationMapRules(updateCustomerRule, customerModel.UpdateRequest{})
	validate.RegisterStructValidationMapRules(deleteCustomerRule, customerModel.DeleteRequest{})
	validate.RegisterStructValidationMapRules(createOrderRule, orderModel.CreateRequest{})
	validate.RegisterStructValidationMapRules(updateOrderRule, orderModel.UpdateRequest{})
	validate.RegisterStructValidationMapRules(deleteOrderRule, orderModel.DeleteRequest{})
	return validate
}

var Validate = NewValidate()

func ValidateLoginRequest(request *loginModel.LoginRequest) error {
	return Validate.Struct(request)
}

func ValidateCreateRequest(request *loginModel.CreateRequest) error {
	return Validate.Struct(request)
}

func ValidateUpdateRequest(request *loginModel.UpdateRequest) error {
	return Validate.Struct(request)
}

func ValidateDeleteRequest(request *loginModel.DeleteRequest) error {
	return Validate.Struct(request)
}

func ValidateCreateCustomerRequest(request *customerModel.CreateRequest) error {
	return Validate.Struct(request)
}

func ValidateUpdateCustomerRequest(request *customerModel.UpdateRequest) error {
	return Validate.Struct(request)
}

func ValidateDeleteCustomerRequest(request *customerModel.DeleteRequest) error {
	return Validate.Struct(request)
}

func ValidateCreateOrderRequest(request *orderModel.CreateRequest) error {
	return Validate.Struct(request)
}

func ValidateUpdateOrderRequest(request *orderModel.UpdateRequest) error {
	return Validate.Struct(request)
}

func ValidateDeleteOrderRequest(request *orderModel.DeleteRequest) error {
	return Validate.Struct(request)
}
