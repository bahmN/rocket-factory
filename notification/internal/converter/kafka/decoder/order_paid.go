package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/bahmN/rocket-factory/notification/internal/model"
	eventsV1 "github.com/bahmN/rocket-factory/shared/pkg/proto/events/v1"
)

type decoderPaid struct{}

func NewOrderPaidDecoder() *decoderPaid {
	return &decoderPaid{}
}

func (d *decoderPaid) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var event model.OrderPaidEvent

	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse event uuid: %w", err)
	}
	event.EventUUID = eventUUID.String()

	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse order uuid: %w", err)
	}
	event.OrderUUID = orderUUID.String()

	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse user uuid: %w", err)
	}
	event.UserUUID = userUUID.String()

	transactionUUID, err := uuid.Parse(pb.TransactionUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse transaction uuid: %w", err)
	}
	event.TransactionUUID = transactionUUID.String()

	event.PaymentMethod = pb.PaymentMethod

	return event, nil
}
