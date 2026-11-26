package horder

import (
	"net/http"

	"github.com/AlexanderTal/order-service/internal/app/entity"
	rhandler "github.com/AlexanderTal/order-service/internal/app/handler/http"
	"github.com/AlexanderTal/order-service/internal/app/service"
	"github.com/AlexanderTal/order-service/internal/pkg/http/httph"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type handler struct {
	srv service.Order
}

func NewHandler(srv service.Order) rhandler.Order {
	return &handler{srv: srv}
}

func (h *handler) Create(c *gin.Context) {
	var req entity.RequestOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Create(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderCreate{
		GUID:          order.GUID,
		UserGUID:      order.UserGUID,
		TotalPrice:    order.TotalPrice,
		DeliveryPrice: order.DeliveryPrice,
		Currency:      order.Currency,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
		Items:         mapItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.GetByGUID(c.Request.Context(), guid)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderGet{
		GUID:          order.GUID,
		UserGUID:      order.UserGUID,
		TotalPrice:    order.TotalPrice,
		DeliveryPrice: order.DeliveryPrice,
		Currency:      order.Currency,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		Items:         mapItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusOK, resp)
}

func (h *handler) Update(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestOrderUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Update(c.Request.Context(), guid, req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderUpdate{
		GUID:          order.GUID,
		UserGUID:      order.UserGUID,
		TotalPrice:    order.TotalPrice,
		DeliveryPrice: order.DeliveryPrice,
		Currency:      order.Currency,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		Items:         mapItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusOK, resp)
}

func (h *handler) Delete(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	if err := h.srv.Delete(c.Request.Context(), guid); err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendEmpty(c.Writer, http.StatusOK)
}

func (h *handler) List(c *gin.Context) {
	var req entity.RequestOrderList

	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
			return
		}
	}

	orders, err := h.srv.List(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	listItems := make([]entity.ResponseOrderListItem, len(orders))
	for i, o := range orders {
		listItems[i] = entity.ResponseOrderListItem{
			GUID:          o.GUID,
			UserGUID:      o.UserGUID,
			TotalPrice:    o.TotalPrice,
			DeliveryPrice: o.DeliveryPrice,
			Currency:      o.Currency,
			Status:        o.Status,
			CreatedAt:     o.CreatedAt,
			UpdatedAt:     o.UpdatedAt,
		}
	}

	resp := entity.ResponseOrderList{Data: listItems}
	httph.SendJSON(c.Writer, http.StatusOK, resp)
}

func mapItems(items []entity.OrderItem) []entity.ResponseOrderItem {
	res := make([]entity.ResponseOrderItem, len(items))
	for i, item := range items {
		res[i] = entity.ResponseOrderItem{
			GUID:        item.GUID,
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}
	return res
}
