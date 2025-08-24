package model

const (
	USER_TYPE_ADMIN = "admin"
	USER_TYPE_USER  = "user"
)

const (
	TRANSACTION_STATUS_PENDING = "pending"
	TRANSACTION_STATUS_SUCCESS = "success"
	TRANSACTION_STATUS_FAILED  = "failed"
)

var ValidTransactionStatuses = []string{
	TRANSACTION_STATUS_PENDING,
	TRANSACTION_STATUS_SUCCESS,
}
