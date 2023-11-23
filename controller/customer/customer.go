package customer

import (
	"fmt"
	"gin-dbo/framework/middleware"
	"gin-dbo/framework/utils"
	mdl "gin-dbo/view/customer"
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

	router.Use(middleware.AuthorizeJWT())
	{
		router.GET("api/customer", u.GetHandler)
		router.GET("api/customer/:id", u.GetByIdHandler)
		router.POST("api/customer", u.CreateHandler)
		router.PUT("api/customer/:id", u.UpdateHandler)
		router.DELETE("api/customer/:id", u.DeleteHandler)
	}
}

// @Summary Get All Customers
// @Description Get All Customers
// @param limit query int false "limit"
// @param page query string false "page"
// @param keyword query string false "name of some customer"
// @Produce json
// @Success 200 {object} mdl.ResponseData
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/customer [get]
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
		c.JSON(err.Code, result)
	}
}

// @Summary Get Customer By Id
// @Description Customer By Id
// @Produce json
// @Success 200 {object} mdl.ResponseDetail
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/customer/{id} [get]
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

// @Summary Create Customer
// @Description Create Some New Customer
// @Accept json
// @Produce json
// @Param request body mdl.CreateRequest true "Sample Create request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/customer [post]
func (u Handler) CreateHandler(c *gin.Context) {
	param := new(mdl.CreateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.createHandler.BadRequest : %v", err.Error())})
		return
	}

	u.logger.Debugf("%+v", param)
	if err := utils.ValidateCreateCustomerRequest(param); err == nil {
		result, err := u.Usecase.Create(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success create data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(err.Code, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.createHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Update Customer
// @Description Update Some Customer
// @Accept json
// @Produce json
// @Param request body mdl.UpdateRequest true "Sample Update request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/customer/{id} [put]
func (u Handler) UpdateHandler(c *gin.Context) {
	var JWT, _ = c.Get(middleware.JwtClaims)
	jwtClaims := JWT.(*middleware.AuthCustomClaims)
	param := new(mdl.UpdateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.updateHandler.BadRequest : %v", err.Error())})
		return
	}

	param.Id = c.Param("id")
	if jwtClaims.Role != "admin" && param.Id != jwtClaims.CustomerId {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.updateHandler.BadRequest : %v", fmt.Errorf("this user can't update this id %s", param.Id))})
		return
	}
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateUpdateCustomerRequest(param); err == nil {
		result, err := u.Usecase.Update(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success update data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(err.Code, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.updateHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Delete Customer
// @Description Delete Some Customer
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/customer/{id} [delete]
func (u Handler) DeleteHandler(c *gin.Context) {
	param := &mdl.DeleteRequest{
		Id: c.Param("id"),
	}
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateDeleteCustomerRequest(param); err == nil {
		result, err := u.Usecase.Delete(c, param)
		if err == nil {
			result.Success = true
			result.Message = "success delete data"
			c.JSON(http.StatusOK, result)
		} else {
			u.logger.Error(err)
			result.Success = false
			result.Message = err.Message.Error()
			c.JSON(err.Code, result)
		}
	} else {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("customer.deleteHandler.BadRequest : %v", err.Error())})
	}
}
