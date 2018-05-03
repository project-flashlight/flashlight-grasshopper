package grasshopper

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PublishResponse struct {
	Status string `json:"status"`
}

type Server struct {
	serv              *http.Server
	geckoboardService *DefGeckoboardService
}

func (me *Server) Start(address string) *httprouter.Router {

	// define routes
	router := httprouter.New()
	router.POST("/publish", me.handleAppStatusPost)

	me.serv = &http.Server{
		Addr:    address,
		Handler: router,
	}

	go func() {
		me.serv.ListenAndServe()
	}()

	return router
}

func (me *Server) Stop() error {
	if me.serv != nil {
		return me.serv.Shutdown(context.Background())
	}
	return nil
}

func (me *Server) handleAppStatusPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	me.geckoboardService.PublishStatus(AppStatus{
		AppName:  "dummy",
		CommitID: "asdasdasdasdas",
		Date:     "2018-03-01",
		Stage:    "Production",
	})
}

func NewServer(geckoboardService DefGeckoboardService) *Server {
	s := Server{
		geckoboardService: &geckoboardService,
	}
	return &s
}
