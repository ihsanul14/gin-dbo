package order

import (
	"fmt"
	mdl "gin-dbo/view/order"

	"github.com/gin-gonic/gin"

	"gin-dbo/controller/customer"
	internal "gin-dbo/framework/error"
	"gin-dbo/framework/utils"
)

type UsecaseModul struct {
	Repo         Repository
	CustomerRepo customer.Repository
}

type Usecase interface {
	Get(ctx *gin.Context, param *mdl.GetRequest) (res mdl.ResponseData, err *internal.Error)
	GetById(ctx *gin.Context, id string) (res mdl.ResponseDetail, err *internal.Error)
	Create(ctx *gin.Context, request *mdl.CreateRequest) (res mdl.GeneralResponse, err *internal.Error)
	Update(ctx *gin.Context, request *mdl.UpdateRequest) (res mdl.GeneralResponse, err *internal.Error)
	Delete(ctx *gin.Context, request *mdl.DeleteRequest) (res mdl.GeneralResponse, err *internal.Error)
}

func NewUsecase(u Repository, c customer.Repository) Usecase {
	return &UsecaseModul{Repo: u, CustomerRepo: c}
}

func (u *UsecaseModul) Get(ctx *gin.Context, param *mdl.GetRequest) (mdl.ResponseData, *internal.Error) {
	var res mdl.ResponseData
	count, err := u.Repo.Count(ctx, param)
	if err != nil {
		return mdl.ResponseData{}, err
	}
	page := utils.GetPage(param.Page)
	totalPage := utils.GetTotalPage(param.Limit, count)

	if page > totalPage {
		return mdl.ResponseData{}, internal.NewError(400, fmt.Errorf("page greater than totalPage"))
	}

	data, err := u.Repo.Get(ctx, param, page)
	if err != nil {
		return mdl.ResponseData{}, err
	}

	res.Data = data
	res.Limit = param.Limit
	res.Page = page
	res.TotalPage = totalPage

	return res, nil
}

func (u *UsecaseModul) GetById(ctx *gin.Context, id string) (mdl.ResponseDetail, *internal.Error) {
	var res mdl.ResponseDetail
	data, err := u.Repo.GetById(ctx, id)
	if err != nil {
		return mdl.ResponseDetail{}, err
	}
	res.Data = data
	return res, nil
}

func (u *UsecaseModul) Create(ctx *gin.Context, param *mdl.CreateRequest) (mdl.GeneralResponse, *internal.Error) {
	var res mdl.GeneralResponse

	_, err := u.CustomerRepo.GetById(ctx, param.CustomerId)
	if err != nil {
		return res, err
	}
	id, err := u.Repo.Create(ctx, param)
	if err != nil {
		return res, err
	}
	res.Id = id
	return res, nil
}

func (u *UsecaseModul) Update(ctx *gin.Context, param *mdl.UpdateRequest) (mdl.GeneralResponse, *internal.Error) {
	var res mdl.GeneralResponse
	_, err := u.CustomerRepo.GetById(ctx, param.CustomerId)
	if err != nil {
		return res, err
	}
	err = u.Repo.Update(ctx, param)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (u *UsecaseModul) Delete(ctx *gin.Context, param *mdl.DeleteRequest) (mdl.GeneralResponse, *internal.Error) {
	var res mdl.GeneralResponse
	err := u.Repo.Delete(ctx, param)
	if err != nil {
		return res, err
	}
	return res, nil
}
