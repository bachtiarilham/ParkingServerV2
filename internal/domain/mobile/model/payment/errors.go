package payment

import "errors"

var ErrPaymentAlreadyCompleted = errors.New("payment already completed")
