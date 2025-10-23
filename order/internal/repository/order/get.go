package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/bahmN/rocket-factory/order/internal/model"
	"github.com/bahmN/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/bahmN/rocket-factory/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.OrderInfo, error) {
	selectBuilder := sq.Select("*").
		From(repoModel.TableOrders).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repoModel.FieldOrderUUID: uuid})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return model.OrderInfo{}, err
	}

	var order repoModel.OrderInfo
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUUIDs,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to get order: %v", err)
		return model.OrderInfo{}, err
	}

	return converter.OrderToModel(order), nil
}
