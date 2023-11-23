package order

import (
	"fmt"
	"gin-dbo/framework/utils"
	models "gin-dbo/model/order"
	view "gin-dbo/view/order"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	internal "gin-dbo/framework/error"
)

type Repo struct {
	Dbconn *gorm.DB
}

type Repository interface {
	Get(ctx *gin.Context, param *view.GetRequest, page int) (res []*models.Order, err *internal.Error)
	Count(ctx *gin.Context, param *view.GetRequest) (res int, err *internal.Error)
	GetById(ctx *gin.Context, id string) (res *models.Order, err *internal.Error)
	Create(ctx *gin.Context, request *view.CreateRequest) (res string, err *internal.Error)
	Update(ctx *gin.Context, request *view.UpdateRequest) (err *internal.Error)
	Delete(ctx *gin.Context, request *view.DeleteRequest) (err *internal.Error)
}

func NewRepository(dbconn *gorm.DB) Repository {
	return &Repo{Dbconn: dbconn}
}

func (r Repo) Get(ctx *gin.Context, param *view.GetRequest, page int) ([]*models.Order, *internal.Error) {
	var (
		res []*models.Order
	)
	query := r.Dbconn
	if param.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+param.Keyword+"%")
	}

	if param.Page > 0 {
		query = query.Offset((page - 1) * param.Limit)
	}

	if param.Limit > 0 {
		query = query.Limit(param.Limit)
	}

	if err := query.Order("created_at desc").Find(&res).Error; err != nil {
		return nil, internal.NewError(500, fmt.Errorf("order.repository.Get : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Count(ctx *gin.Context, param *view.GetRequest) (int, *internal.Error) {
	var (
		res int
	)
	query := r.Dbconn.Select("COUNT(1) as total").Model(&models.Order{})
	if param.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+param.Keyword+"%")
	}

	if err := query.Pluck("total", &res).Error; err != nil {
		return 0, internal.NewError(500, fmt.Errorf("customer.repository.Count : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) GetById(ctx *gin.Context, id string) (*models.Order, *internal.Error) {
	var (
		res *models.Order
		err error
	)
	query := r.Dbconn.Find(&res).Where("id = ?", id)
	if err = query.Error; err != nil {
		return nil, internal.NewError(500, fmt.Errorf("order.repository.GetById : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Create(ctx *gin.Context, param *view.CreateRequest) (string, *internal.Error) {
	var err error

	uid := uuid.New().String()
	now := utils.FormatTime()
	query := r.Dbconn.Create(models.Order{Id: uid, CustomerId: param.CustomerId, Name: param.Name, Qty: param.Qty, CreatedAt: now, UpdatedAt: now})
	if err = query.Error; err != nil {
		return "", internal.NewError(500, fmt.Errorf("order.repository.Create : %v", err.Error()))
	}

	return uid, nil
}

func (r Repo) Update(ctx *gin.Context, param *view.UpdateRequest) *internal.Error {
	err := r.Dbconn.Updates(models.Order{Id: param.Id, CustomerId: param.CustomerId, Name: param.Name, Qty: param.Qty, UpdatedAt: utils.FormatTime()}).Error
	if err != nil {
		return internal.NewError(500, fmt.Errorf("order.repository.Update : %v", err.Error()))
	}
	return nil
}

func (r Repo) Delete(ctx *gin.Context, param *view.DeleteRequest) *internal.Error {
	err := r.Dbconn.Delete(models.Order{Id: param.Id}).Error
	if err != nil {
		return internal.NewError(500, fmt.Errorf("order.repository.Delete : %v", err.Error()))
	}
	return nil
}
