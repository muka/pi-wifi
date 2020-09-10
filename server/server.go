package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//NewHTTPServer init an http server wrapper
func NewHTTPServer(
	onConnect func(string) (string, error),
	onStatus func() (string, error),
	onListAP func() ([]AccessPoint, error),
) HTTPServer {
	return HTTPServer{
		onConnect: onConnect,
		onStatus:  onStatus,
		onListAP:  onListAP,
	}
}

// HTTPServer wraps http server API
type HTTPServer struct {
	onConnect func(string) (string, error)
	onListAP  func() ([]AccessPoint, error)
	onStatus  func() (string, error)
}

// Serve starts an http server
func (srv *HTTPServer) Serve() error {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed to read body")
			w.WriteHeader(500)
			return
		}

		if len(body) == 0 {
			log.Errorf("Empty request body")
			w.WriteHeader(400)
			return
		}

		msg := connectRequest{}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Errorf("Failed to parse JSON")
			w.WriteHeader(400)
			return
		}

		res, err := srv.onConnect(msg.Payload)
		if err != nil {
			log.Errorf("Failed to connect: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "status": "%s" }`, res)
	})

	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(200)

		status, err := srv.onStatus()
		if err != nil {
			log.Errorf("failed to get status: %s", err)
			w.WriteHeader(500)
			return
		}

		res := statusResponse{
			Status: status,
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Errorf("failed to encode response: %s", err)
			w.WriteHeader(500)
			return
		}

	})

	router.HandleFunc("/listap", func(w http.ResponseWriter, r *http.Request) {

		list, err := srv.onListAP()
		if err != nil {
			log.Errorf("failed to get status: %s", err)
			w.WriteHeader(500)
			return
		}

		res := listAPResponse{
			AccessPoints: list,
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Errorf("failed to encode response: %s", err)
			w.WriteHeader(500)
			return
		}

	})

	port := fmt.Sprintf(":%d", viper.GetInt("http_port"))
	log.Debugf("Serving on %s", port)
	return http.ListenAndServe(port, router)
}
