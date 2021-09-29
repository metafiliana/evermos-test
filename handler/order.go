package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metafiliana/evermos-test/order"
	"github.com/metafiliana/evermos-test/response"
)

type handler struct {
	orderSvc order.Service
}

type Handler interface {
	CheckoutItems(c *gin.Context)
	CreateOrder(c *gin.Context)
}

func NewHandler(orderSvc order.Service) Handler {
	return &handler{orderSvc: orderSvc}
}

func (h *handler) CheckoutItems(c *gin.Context) {
	var req order.CheckoutOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		response.SendResponseWithError(c, err, `error bad request CheckoutItems`)
		return
	}

	err := h.orderSvc.CheckoutItems(&req)
	if err != nil {
		response.SendResponseWithError(c, err, `error failed to request CheckoutItems`)
		return
	}
	response.SendResponse(c, http.StatusOK, `SUCCESS`, nil)
}

func (h *handler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		response.SendResponseWithError(c, err, `error bad request CreateOrder`)
		return
	}

	err := h.orderSvc.CreateOrder(&req)
	if err != nil {
		response.SendResponseWithError(c, err, `error failed to request CreateOrder`)
		return
	}
	response.SendResponse(c, http.StatusOK, `SUCCESS`, nil)
}
