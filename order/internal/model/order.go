package model

const (
	OrderStatusPENDINGPAYMENT string = "PENDING_PAYMENT"
	OrderStatusPAID           string = "PAID"
	OrderStatusCANCELLED      string = "CANCELLED"
	OrderStatusASSEMBLED      string = "ASSEMBLED"
)

type CreateOrderReq struct {
	UserUUID  string
	PartsUUID []string
}

type CreateOrderResp struct {
	OrderUUID  string
	TotalPrice float64
}

type OrderInfo struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
}
