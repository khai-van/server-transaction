package main

import (
	"app/api"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"testing"

	// "app/transaction"

	"github.com/kataras/iris/v12/httptest"
)

const (
	numberRequest = 10000        // number of register user
	mainAccount   = "nVrsWDeiTX" // default user with big balance
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestApp(t *testing.T) {
	const defaultConfigFilename = "server.yml"
	var wg sync.WaitGroup

	var serverConfig api.Configuration
	serverConfig.BindFile(defaultConfigFilename)
	app := api.NewServer(serverConfig)

	e := httptest.New(t, app.Application)
	// get first balance
	resp := e.GET("/detail").WithQuery("account", mainAccount).Expect()
	resp.Status(httptest.StatusOK)
	balance := uint64(resp.JSON().Object().Raw()["payload"].(map[string]interface{})["user"].(map[string]interface{})["balance"].(float64))
	fmt.Println(balance)

	for i := 0; i < numberRequest; i++ {
		wg.Add(1)
		account := RandStringRunes(10)
		go func() { // send request in diffenrence goroutine no waitting
			e.GET("/register").WithQuery("account", account).Expect().Status(httptest.StatusOK) // test register
			e.GET("/detail").WithQuery("account", account).Expect().Status(httptest.StatusOK)   // test get detail
			// test transfer main account to this account
			e.GET("/transfer").
				WithQuery("from", mainAccount).
				WithQuery("to", account).
				WithQuery("amount", 100).
				Expect().
				Status(httptest.StatusOK)
			// test transfer this account back to main account
			e.GET("/transfer").
				WithQuery("from", account).
				WithQuery("to", mainAccount).
				WithQuery("amount", 50).
				Expect().
				Status(httptest.StatusOK)

			wg.Done()
		}()

	}

	wg.Wait()
	// get detail main account
	newResp := e.GET("/detail").WithQuery("account", mainAccount).Expect()
	newResp.Status(httptest.StatusOK)
	newBalance := uint64(newResp.JSON().Object().Raw()["payload"].(map[string]interface{})["user"].(map[string]interface{})["balance"].(float64))
	// recalculate main account and compare with in system db
	if (balance - (numberRequest*100 - numberRequest*50)) != newBalance {
		t.Error("Wrong expect Balance ", (balance - (numberRequest*100 - numberRequest*50)), newBalance)
	}
}
