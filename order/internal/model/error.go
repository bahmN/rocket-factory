package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")

var ErrPartsNotFound = errors.New("parts not found")

var ErrOrderPaidOrCanceled = errors.New("order paid or canceled")

var ErrEmptyUUID = errors.New("empty UUID")
