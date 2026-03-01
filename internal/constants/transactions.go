package constants

const (
	TxStatusInProgress int16 = 0
	TxStatusApproved   int16 = 1
	TxStatusCancelled  int16 = 2
	TxStatusFailed     int16 = 3
)

const (
	PaymentTypeDeposit int16 = 1
	PaymentTypeCharge  int16 = 2
	PaymentTypeRefund  int16 = 3
	PaymentTypeBonus   int16 = 4
)

const (
	PaymentMethodCard int16 = 1
	PaymentMethodBank int16 = 2
)
