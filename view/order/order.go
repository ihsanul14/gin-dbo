package order

import "gin-dbo/model/order"

type GetRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}
type CreateRequest struct {
	CustomerId string `json:"customerId"`
	Name       string `json:"name"`
	Qty        int64  `json:"qty"`
}

type UpdateRequest struct {
	Id         string `json:"id" swaggerignore:"true"`
	CustomerId string `json:"customerId"`
	Name       string `json:"name"`
	Qty        int64  `json:"qty"`
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
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *order.Order `json:"data"`
}

type ResponseData struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	Data      []*order.Order `json:"data"`
	Limit     int            `json:"limit"`
	Page      int            `json:"page"`
	TotalPage int            `json:"totalPage"`
}
