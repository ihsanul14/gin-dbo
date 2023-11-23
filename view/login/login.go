package login

import "gin-dbo/model/login"

type GetRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}
type CreateRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	CustomerId string `json:"customerId" swaggerignore:"true"`
}

type UpdateRequest struct {
	Username   string `json:"username" swaggerignore:"true"`
	Password   string `json:"password"`
	CustomerId string `json:"customerId,omitempty"`
	Role       string `json:"role"`
}

type DeleteRequest struct {
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GeneralResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response400 struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"invalid request"`
}

type Response500 struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"something went wrong"`
}
type ResponseLogin struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token,omitempty"`
	} `json:"data,omitempty"`
}

type ResponseDetail struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *login.User `json:"data"`
}

type ResponseData struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Data      []*login.User `json:"data"`
	Limit     int           `json:"limit"`
	Page      int           `json:"page"`
	TotalPage int           `json:"totalPage"`
}
