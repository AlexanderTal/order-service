package porder

import (
	"context"
	"errors"
	"time"

	"github.com/AlexanderTal/order-service/internal/app/entity"
	"github.com/AlexanderTal/order-service/internal/app/repository"
	rcpostgres "github.com/AlexanderTal/order-service/internal/app/repository/conn/postgres"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type repoPg struct {
	conn *rcpostgres.Client
}

func NewRepo(client *rcpostgres.Client) repository.Order {
	return &repoPg{conn: client}
}

func (r *repoPg) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.conn.InsideTx(ctx, fn)
}

func (r *repoPg) Create(ctx context.Context, order entity.Order) error {
	return r.conn.GetDB(ctx).WithContext(ctx).Create(&order).Error
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	var order entity.Order
	err := r.conn.GetDB(ctx).WithContext(ctx).
		Preload("Items").
		Where("guid = ?", guid.String()).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Order{}, entity.ErrNotFound
		}
		return entity.Order{}, err
	}
	return order, nil
}

func (r *repoPg) Update(ctx context.Context, order entity.Order) error {
	res := r.conn.GetDB(ctx).WithContext(ctx).
		Model(&entity.Order{}).
		Where("guid = ?", order.GUID.String()).
		Updates(map[string]any{
			"status":     order.Status,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	res := r.conn.GetDB(ctx).WithContext(ctx).
		Where("guid = ?", guid.String()).
		Delete(&entity.Order{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (r *repoPg) List(ctx context.Context, status *string, userGUID *uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	query := r.conn.GetDB(ctx).WithContext(ctx).Model(&entity.Order{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userGUID != nil {
		query = query.Where("user_guid = ?", userGUID.String())
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
