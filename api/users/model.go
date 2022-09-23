package users

type UserRequest struct {
	Account string `json:"account" validate:"required"`
}

type TransferRequest struct {
	From   string `json:"from" validate:"required"`
	To     string `json:"to" validate:"required"`
	Amount uint64 `json:"amount" validate:"required"`
}

type Response struct {
	Message string      `json:"msg"`
	Payload interface{} `json:"payload"`
}
