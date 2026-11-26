package sorder

import (
	"context"
	"time"

	"github.com/AlexanderTal/order-service/internal/app/entity"
	"github.com/AlexanderTal/order-service/internal/app/repository"
	"github.com/AlexanderTal/order-service/internal/app/service"
	"github.com/gofrs/uuid"
)

type srv struct {
	repoOrder repository.Order
}

func NewService(repoOrder repository.Order) service.Order {
	return &srv{repoOrder: repoOrder}
}

func (s *srv) Create(ctx context.Context, req entity.RequestOrderCreate) (entity.Order, error) {
	now := time.Now()
	var totalPrice int64
	items := make([]entity.OrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		totalPrice += int64(itemReq.Quantity) * itemReq.UnitPrice

		itemGUID := uuid.Must(uuid.NewV4())
		items = append(items, entity.OrderItem{
			GUID:        itemGUID,
			ProductGUID: itemReq.ProductGUID,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	order := entity.Order{
		GUID:          uuid.Must(uuid.NewV4()),
		UserGUID:      req.UserGUID,
		TotalPrice:    totalPrice,
		DeliveryPrice: 0,
		Currency:      req.Currency,
		Status:        "pending",
		CreatedAt:     now,
		UpdatedAt:     now,
		Items:         items,
	}

	if err := s.repoOrder.Create(ctx, order); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *srv) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	return s.repoOrder.GetByGUID(ctx, guid)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestOrderUpdate) (entity.Order, error) {
	var updatedOrder entity.Order

	err := s.repoOrder.InsideTx(ctx, func(txCtx context.Context) error {
		order, err := s.repoOrder.GetByGUID(txCtx, guid)
		if err != nil {
			return err
		}

		order.Status = req.Status
		order.UpdatedAt = time.Now()

		if err := s.repoOrder.Update(txCtx, order); err != nil {
			return err
		}

		updatedOrder = order
		return nil
	})
	if err != nil {
		return entity.Order{}, err
	}

	return updatedOrder, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoOrder.InsideTx(ctx, func(txCtx context.Context) error {
		_, err := s.repoOrder.GetByGUID(txCtx, guid)
		if err != nil {
			return err
		}
		return s.repoOrder.Delete(txCtx, guid)
	})
}

func (s *srv) List(ctx context.Context, req entity.RequestOrderList) ([]entity.Order, error) {
	return s.repoOrder.List(ctx, req.Status, req.UserGUID)
}
