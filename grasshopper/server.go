package grasshopper

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
}

func (me *Server) Start(address string) *httprouter.Router {

	// define routes
	router := httprouter.New()
	router.GET("/", me.handleWelcome)
	router.POST("/publish", me.handleAppStatusPost)

	go func() {
		http.ListenAndServe(address, router)
	}()

	return router
}

func (me *Server) handleAppStatusPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("Received something")
	w.Write([]byte("Handle App Status POST"))
}

func (me *Server) handleWelcome(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("Received something")
	fmt.Fprint(w, "Welcome!\n")
}

func NewServer() *Server {
	s := Server{}
	return &s
}
