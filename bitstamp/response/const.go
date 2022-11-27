package response

// Transaction types
const (
	TransactionTypeDeposit            = 0
	TransactionTypeWithdrawal         = 1
	TransactionTypeMarketTrade        = 2
	TransactionTypeSubAccountTransfer = 14
)

// Order types
const (
	OrderTypeBuy  = 0
	OrderTypeSell = 1
)

// Withdrawal types
const (
	WithdrawalTypeSEPA = 0
	WithdrawalTypeBTC  = 1
	WithdrawalTypeWIRE = 2
	WithdrawalTypeXRP  = 14
	WithdrawalTypeLTC  = 15
	WithdrawalTypeETH  = 15
)
