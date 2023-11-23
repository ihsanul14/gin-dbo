package order

import (
	"fmt"
	"gin-dbo/framework/middleware"
	"gin-dbo/framework/utils"
	mdl "gin-dbo/view/order"
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
		router.GET("api/order", u.GetHandler)
		router.GET("api/order/:id", u.GetByIdHandler)
		router.POST("api/order", u.CreateHandler)
		router.PUT("api/order/:id", u.UpdateHandler)
		router.DELETE("api/order/:id", u.DeleteHandler)
	}
}

// @Summary Get All Orders
// @Description Get All Orders
// @Produce json
// @Success 200 {object} mdl.ResponseData
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/order [get]
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

// @Summary Get Order By Id
// @Description Get Order By Id
// @Produce json
// @Success 200 {object} mdl.ResponseData
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/order/{id} [get]
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

// @Summary Create Order
// @Description Create Some New Orders
// @Accept json
// @Produce json
// @Param request body mdl.CreateRequest true "Sample Create request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/order [post]
func (u Handler) CreateHandler(c *gin.Context) {
	param := new(mdl.CreateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("order.createHandler.BadRequest : %v", err.Error())})
		return
	}

	u.logger.Debugf("%+v", param)
	if err := utils.ValidateCreateOrderRequest(param); err == nil {
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
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("order.createHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Update Order
// @Description Update Some Orders
// @Accept json
// @Produce json
// @Param request body mdl.UpdateRequest true "Sample Update request payload"
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/order/{id} [put]
func (u Handler) UpdateHandler(c *gin.Context) {
	param := new(mdl.UpdateRequest)
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("order.updateHandler.BadRequest : %v", err.Error())})
		return
	}
	param.Id = c.Param("id")
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateUpdateOrderRequest(param); err == nil {
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
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("order.updateHandler.BadRequest : %v", err.Error())})
	}
}

// @Summary Delete Order
// @Description Delete Some Orders
// @Accept json
// @Produce json
// @Security jwt
// @Success 200 {object} mdl.GeneralResponse
// @Failure 400 {object} mdl.Response400
// @Failure 401 {object} middleware.Response
// @Failure 500 {object} mdl.Response500
// @Router /api/order/{id} [delete]
func (u Handler) DeleteHandler(c *gin.Context) {
	param := &mdl.DeleteRequest{
		Id: c.Param("id"),
	}
	u.logger.Debugf("%+v", param)

	if err := utils.ValidateDeleteOrderRequest(param); err == nil {
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
		c.JSON(http.StatusBadRequest, mdl.GeneralResponse{Success: false, Message: fmt.Sprintf("order.deleteHandler.BadRequest : %v", err.Error())})
	}
}
