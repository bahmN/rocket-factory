package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bahmN/rocket-factory/order/internal/model"
	repoModel "github.com/bahmN/rocket-factory/order/internal/repository/model"
	"github.com/samber/lo"
)

func (r *repository) Update(ctx context.Context, uuid string, info model.OrderInfo) error {
	updateBuilder := sq.Update(repoModel.TableOrders).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repoModel.FieldOrderUUID: uuid}).
		Set(repoModel.FieldUserUUID, info.UserUUID).
		Set(repoModel.FieldPartUUIDs, info.PartUUIDs).
		Set(repoModel.FieldTotalPrice, info.TotalPrice).
		Set(repoModel.FieldStatus, info.Status).
		Set(repoModel.FieldUpdatedAt, time.Now())

	if lo.IsNotEmpty(info.TransactionUUID) {
		updateBuilder = updateBuilder.Set(repoModel.FieldTransactionUUID, info.TransactionUUID)
	}

	if lo.IsNotEmpty(info.PaymentMethod) {
		updateBuilder = updateBuilder.Set(repoModel.FieldPaymentMethod, info.PaymentMethod)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build sql query: %v", err)
		return err
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	log.Printf("order with UUID %v updated successfully", uuid)
	return nil
}
