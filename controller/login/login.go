package login

import (
	"errors"
	"fmt"
	"gin-dbo/framework/middleware"
	"gin-dbo/framework/utils"
	mdl "gin-dbo/view/login"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Usecase Usecase
	logger  *logrus.Logger
}

// @SecurityDefinitions jwt
func Router(router *gin.Engine, uc Usecase, logger *logrus.Logger) {
	u := Handler{Usecase: uc, logger: logger}

	router.POST("api/register", u.CreateHandler)
	router.POST("api/login", u.LoginHandler)
	router.Use(middleware.AuthorizeJWT())
	{
		router.GET("api/user", u.GetHandler)
		router.GET("api/user/:id", u.GetByIdHandler)
		router.POST("api/user", u.CreateHandler)
		router.PUT("api/user/:id", u.UpdateHandler)
		router.DELETE("api/user/:id", u.DeleteHandler)
	}
}

// @Summary Register
// @Description Create Some New Users
// @Accept json
// @Produce json
// @Param request body mdl.CreateRequest true "Sample Create request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/register [post]
func (u Handler) RegisterHandler(c *gin.Context) {
	param := new(mdl.CreateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("user.registerHandler.BadRequest : %v", err.Error())})
		return
	}

	u.logger.Debugf("%+v", param)
	result, err := u.Usecase.Create(c, param)
	if err == nil {
		result.Success = true
		result.Message = "success create data"
		c.JSON(http.StatusOK, result)
	} else {
		u.logger.Error(err)
		result.Success = false
		result.Message = err.Message.Error()
		c.JSON(http.StatusInternalServerError, result)
	}
}

// @Summary Login
// @Description Handle Login of Some Users
// @Accept json
// @Produce json
// @Param request body mdl.LoginRequest true "Sample Login request payload"
// @Success 200 {object} mdl.ResponseLogin
// @Failure 400 {object} mdl.GeneralResponse
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.GeneralResponse
// @Router /api/login [post]
func (u Handler) LoginHandler(c *gin.Context) {
	param := new(mdl.LoginRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("user.loginHandler.BadRequest : %v", err.Error())})
		return
	}
	u.logger.Debug(param)

	if utils.ValidateLoginRequest(param) == nil {
		result, err := u.Usecase.Login(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success login"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(http.StatusInternalServerError, result)
		}
	} else {
		errn := errors.New("mandatory field is missing")
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("users.delivery.createHandler.BadRequest : %v", errn)})
	}
}

// @Summary Get All Users
// @Description Get All Users
// @Produce json
// @Success 200 {object} mdl.ResponseData
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/user [get]
func (u Handler) GetHandler(c *gin.Context) {
	limit, err := utils.GetLimit(c.Query(utils.Limit))
	if err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.getHandler.BadRequest : %v", err.Message.Error())})
		return
	}

	page, err := utils.GetTargetPage(c.Query(utils.Page))
	if err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.getHandler.BadRequest : %v", err.Message.Error())})
		return
	}

	param := &mdl.GetRequest{
		Keyword: c.Query(utils.Keyword),
		Limit:   limit,
		Page:    page,
	}

	result, err := u.Usecase.Get(c, param)
	if err == nil {
		result.Success = true
		result.Message = "success retrieve data"
		c.JSON(http.StatusOK, result)
	} else {
		u.logger.Error(err)
		result.Success = false
		result.Message = err.Message.Error()
		c.JSON(http.StatusInternalServerError, result)
	}
}

// @Summary Get User By Id
// @Description Get User By Id
// @Produce json
// @Success 200 {object} mdl.ResponseDetail
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/user/{id} [get]
func (u Handler) GetByIdHandler(c *gin.Context) {
	id := c.Param("id")
	result, err := u.Usecase.GetById(c, id)
	if err == nil {
		result.Success = true
		result.Message = "success retrieve data"
		c.JSON(http.StatusOK, result)
	} else {
		u.logger.Error(err)
		result.Success = false
		result.Message = err.Message.Error()
		c.JSON(err.Code, result)
	}
}

// @Summary Create User
// @Description Create Some New Users
// @Accept json
// @Produce json
// @Param request body mdl.CreateRequest true "Sample Create request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/user [post]
func (u Handler) CreateHandler(c *gin.Context) {
	param := new(mdl.CreateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("user.createHandler.BadRequest : %v", err.Error())})
		return
	}

	u.logger.Debugf("%+v", param)
	if err := utils.ValidateCreateRequest(param); err == nil {
		result, err := u.Usecase.Create(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success create data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(http.StatusInternalServerError, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("users.delivery.createHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Update User
// @Description Update Some Users
// @Accept json
// @Produce json
// @Param request body mdl.UpdateRequest true "Sample Update request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/user/{id} [put]
func (u Handler) UpdateHandler(c *gin.Context) {
	param := new(mdl.UpdateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("user.updateHandler.BadRequest : %v", err.Error())})
		return
	}
	param.Username = c.Param("id")
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateUpdateRequest(param); err == nil {
		result, err := u.Usecase.Update(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success update data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(http.StatusInternalServerError, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("users.delivery.updateHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Delete User
// @Description Delete Some Users
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/user/{id} [delete]
func (u Handler) DeleteHandler(c *gin.Context) {
	param := &mdl.DeleteRequest{
		Username: c.Param("id"),
	}
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateDeleteRequest(param); err == nil {
		result, err := u.Usecase.Delete(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success delete data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(http.StatusInternalServerError, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("users.delivery.deleteHandler.BadRequest : %v", err.Error())})
	}
}
