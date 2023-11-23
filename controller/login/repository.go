package login

import (
	"fmt"
	"gin-dbo/framework/middleware"
	"gin-dbo/framework/utils"
	models "gin-dbo/model/login"
	view "gin-dbo/view/login"

	internal "gin-dbo/framework/error"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repo struct {
	Dbconn *gorm.DB
}

type Repository interface {
	Login(ctx *gin.Context, request *view.LoginRequest) (res view.ResponseLogin, err *internal.Error)
	Get(ctx *gin.Context, request *view.GetRequest, page int) (res []*models.User, err *internal.Error)
	Count(ctx *gin.Context, request *view.GetRequest) (res int, err *internal.Error)
	GetById(ctx *gin.Context, id string) (res *models.User, err *internal.Error)
	Create(ctx *gin.Context, request *view.CreateRequest) (res view.GeneralResponse, err *internal.Error)
	Update(ctx *gin.Context, request *view.UpdateRequest) (res view.GeneralResponse, err *internal.Error)
	Delete(ctx *gin.Context, request *view.DeleteRequest) (res view.GeneralResponse, err *internal.Error)
}

func NewRepository(dbconn *gorm.DB) Repository {
	return &Repo{Dbconn: dbconn}
}

func (r Repo) Login(ctx *gin.Context, param *view.LoginRequest) (view.ResponseLogin, *internal.Error) {
	var (
		result *models.User
		res    view.ResponseLogin
		err    error
	)
	query := r.Dbconn.Model(&models.User{}).Where("username = ? and password = ?", param.Username, utils.Encrypt(param.Password))
	err = query.Scan(&result).Error
	if err != nil {
		return res, internal.NewError(500, fmt.Errorf("login.repository.GetDetail : %v", err.Error()))
	}

	if result == nil {
		return res, internal.NewError(400, fmt.Errorf("username or password invalid"))
	}

	tokenString := middleware.JWTAuthService().GenerateToken(result)
	res.Data.Token = tokenString
	return res, nil
}

func (r Repo) Get(ctx *gin.Context, param *view.GetRequest, page int) ([]*models.User, *internal.Error) {
	var (
		res []*models.User
	)
	query := r.Dbconn.Select("username, role, customer_id, created_at, updated_at")
	if param.Keyword != "" {
		query = query.Where("username LIKE ?", "%"+param.Keyword+"%")
	}

	if param.Page > 0 {
		query = query.Offset((page - 1) * param.Limit)
	}

	if param.Limit > 0 {
		query = query.Limit(param.Limit)
	}

	if err := query.Order("created_at desc").Find(&res).Error; err != nil {
		return nil, internal.NewError(500, fmt.Errorf("user.repository.Get : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Count(ctx *gin.Context, param *view.GetRequest) (int, *internal.Error) {
	var (
		res int
	)
	query := r.Dbconn.Select("COUNT(1) as total").Model(&models.User{})
	if param.Keyword != "" {
		query = query.Where("username LIKE ?", "%"+param.Keyword+"%")
	}

	if err := query.Pluck("total", &res).Error; err != nil {
		return 0, internal.NewError(500, fmt.Errorf("user.repository.Count : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) GetById(ctx *gin.Context, id string) (*models.User, *internal.Error) {
	var (
		res *models.User
		err error
	)
	query := r.Dbconn.Find(&res).Where("username = ?", id)
	if err = query.Error; err != nil {
		return res, internal.NewError(500, fmt.Errorf("login.repository.GetById : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Create(ctx *gin.Context, param *view.CreateRequest) (view.GeneralResponse, *internal.Error) {
	var (
		res view.GeneralResponse
		err error
	)
	now := utils.FormatTime()
	query := r.Dbconn.Create(models.User{Username: param.Username, Password: utils.Encrypt(param.Password), Role: param.Role, CustomerId: param.CustomerId, CreatedAt: now, UpdatedAt: now})
	if err = query.Error; err != nil {
		return res, internal.NewError(500, fmt.Errorf("login.repository.Create : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Update(ctx *gin.Context, param *view.UpdateRequest) (view.GeneralResponse, *internal.Error) {
	var res view.GeneralResponse
	err := r.Dbconn.Updates(models.User{Username: param.Username, Password: utils.Encrypt(param.Password), Role: param.Role, UpdatedAt: utils.FormatTime()}).Error
	if err != nil {
		return res, internal.NewError(500, fmt.Errorf("login.repository.Update : %v", err.Error()))
	}
	return res, nil
}

func (r Repo) Delete(ctx *gin.Context, param *view.DeleteRequest) (view.GeneralResponse, *internal.Error) {
	var res view.GeneralResponse
	err := r.Dbconn.Delete(models.User{Username: param.Username}).Error
	if err != nil {
		return res, internal.NewError(500, fmt.Errorf("login.repository.Delete : %v", err.Error()))
	}
	return res, nil
}
