package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/vwdilab/flashlight-grasshopper/grasshopper"
)

var appCtx context.Context
var router *httprouter.Router
var wg sync.WaitGroup

func init() {
	appCtx = context.Background()
	wg = sync.WaitGroup{}
	wg.Add(1)
}

func main() {

	server := grasshopper.NewServer()

	p, ok := loadPort()
	if !ok {
		fmt.Println("Failed to start server: PORT not set in environment")
		return
	}

	router = server.Start(":" + p)

	fmt.Printf("Server started in port %s\n", p)

	wg.Done()

	for {
		select {
		case <-appCtx.Done():
			return
		}
	}
}

func loadPort() (string, bool) {
	p := os.Getenv("PORT")
	return p, p != ""
}
