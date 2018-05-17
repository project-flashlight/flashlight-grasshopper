package grasshopper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PublishResponse struct {
	Status string `json:"status"`
}

type Server struct {
	serv              *http.Server
	geckoboardService GeckoboardService
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
	res := AppStatus{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := me.geckoboardService.PublishStatus(res)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				panic("There was an error parsing the response error body")
			}

			bodyString := string(bodyBytes)

			responseMsg := fmt.Sprintf("Failed to invoke geckoboard: %v", bodyString)
			http.Error(w, responseMsg, resp.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

}

func NewServer(geckoboardService GeckoboardService) *Server {
	s := Server{
		geckoboardService: geckoboardService,
	}
	return &s
}
