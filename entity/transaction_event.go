package entity

type (
	TransactionEvent struct {
		ID     string
		Type   TransactionType
		Params interface{}
	}

	TopUp struct {
		UserID uint
		Amount uint64
	}

	TransactionType string
)

const (
	TransactionTypeTopUp TransactionType = "top_up"
)
