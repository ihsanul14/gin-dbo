package customer

import "gin-dbo/model/customer"

type GetRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}
type CreateRequest struct {
	Name string `json:"name"`
}

type UpdateRequest struct {
	Id   string `json:"id" swaggerignore:"true"`
	Name string `json:"name"`
}

type DeleteRequest struct {
	Id string `json:"id"`
}

type GeneralResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Id      string `json:"id,omitempty"`
}

type Response400 struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"invalid request"`
}

type Response500 struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"something went wrong"`
}

type ResponseDetail struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *customer.Customer `json:"data"`
}

type ResponseData struct {
	Success   bool                 `json:"success"`
	Message   string               `json:"message"`
	Data      []*customer.Customer `json:"data"`
	Limit     int                  `json:"limit"`
	Page      int                  `json:"page"`
	TotalPage int                  `json:"totalPage"`
}
