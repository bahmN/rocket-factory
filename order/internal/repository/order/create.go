package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/bahmN/rocket-factory/order/internal/model"
	repoModel "github.com/bahmN/rocket-factory/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.OrderInfo) error {
	insertBuilder := sq.Insert(repoModel.TableOrders).
		PlaceholderFormat(sq.Dollar).
		Columns(
			repoModel.FieldOrderUUID,
			repoModel.FieldUserUUID,
			repoModel.FieldPartUUIDs,
			repoModel.FieldTotalPrice,
			repoModel.FieldStatus).
		Values(
			info.OrderUUID,
			info.UserUUID,
			info.PartUUIDs,
			info.TotalPrice,
			info.Status).
		Suffix("RETURNING " + repoModel.FieldOrderUUID)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&info.OrderUUID)
	if err != nil {
		log.Printf("failed to insert order: %v", err)
		return err
	}

	return nil
}
