package users

import (
	"app/transaction"
	"errors"

	"github.com/kataras/iris/v12"
)

type UserService struct {
	Transaction transaction.Repository
}

func (service *UserService) Configure(r iris.Party) {
	// register account [GET] :/register?account={string}
	r.Get("/register", service.userRegister)
	// get detail account [GET] :/detail?account={string}
	r.Get("/detail", service.getDetail)
	// make a transfer [GET] :/register?from={string}&to={string}&ammount={uint64}
	r.Get("/transfer", service.transfer)
}

func (service *UserService) userRegister(ctx iris.Context) {
	// parse request
	var req UserRequest
	ctx.ReadQuery(&req)
	if len(req.Account) == 0 {
		ctx.StopWithError(500, errors.New("missing fields account or empty"))

		return
	}

	// call func create user
	res, err := service.Transaction.CreateUser(req.Account)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	// response
	ctx.JSON(Response{
		Message: "Success",
		Payload: res,
	})
}

func (service *UserService) getDetail(ctx iris.Context) {
	// parse request
	var req UserRequest
	ctx.ReadQuery(&req)
	if len(req.Account) == 0 {
		ctx.StopWithError(500, errors.New("missing fields account or empty"))

		return
	}
	// get user
	resUser, err := service.Transaction.GetUser(req.Account)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	// get history transaction
	resTransaction, err := service.Transaction.GetListTransaction(req.Account)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	//response
	ctx.JSON(Response{
		Message: "Success",
		Payload: iris.Map{
			"user":         resUser,
			"transactions": resTransaction,
		},
	})
}

func (service *UserService) transfer(ctx iris.Context) {
	// parse request
	var req TransferRequest
	ctx.ReadQuery(&req)
	if len(req.From) == 0 || len(req.To) == 0 || req.From == req.To || req.Amount == 0 {
		ctx.StopWithError(500, errors.New("missing fields or empty"))
		return
	}
	// make a transfer
	res, err := service.Transaction.Transfer(req.From, req.To, req.Amount)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	// response
	ctx.JSON(Response{
		Message: "Success",
		Payload: res,
	})
}
