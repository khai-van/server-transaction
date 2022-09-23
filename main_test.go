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
	numberRequest = 100
	mainAccount = "nVrsWDeiTX"
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

	resp := e.GET("/detail").WithQuery("account", mainAccount).Expect()
	resp.Status(httptest.StatusOK)
	balance := uint64(resp.JSON().Object().Raw()["payload"].(map[string]interface{})["user"].(map[string]interface{})["balance"].(float64))
	fmt.Println(balance)

	
	for i:=0; i<numberRequest; i++ {
		wg.Add(1)
		account := RandStringRunes(10)
		go func() {
			e.GET("/register").WithQuery("account", account).Expect().Status(httptest.StatusOK)
			e.GET("/detail").WithQuery("account", account).Expect().Status(httptest.StatusOK)
			
			e.GET("/transfer").
				WithQuery("from", mainAccount).
				WithQuery("to", account).
				WithQuery("amount", 100).
				Expect().
				Status(httptest.StatusOK)
			
			e.GET("/transfer").
				WithQuery("from", account).
				WithQuery("to", mainAccount).
				WithQuery("amount", 50).
				Expect().
				Status(httptest.StatusOK)			

			wg.Done()
			} ()
		
		
	}

	wg.Wait()

	newResp := e.GET("/detail").WithQuery("account", mainAccount).Expect()
	newResp.Status(httptest.StatusOK)
	newBalance := uint64(newResp.JSON().Object().Raw()["payload"].(map[string]interface{})["user"].(map[string]interface{})["balance"].(float64))

	if (balance - (numberRequest * 100 - numberRequest * 50)) != newBalance {
		t.Errorf("Wrong expect Balance")
	}
}
