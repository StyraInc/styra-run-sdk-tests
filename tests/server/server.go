package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const host = "127.0.0.1"

type Gateway struct {
	Url string `json:"gateway_url"`
}

type Settings struct {
	Port int
}

type Server interface {
	Listen() error
	Shutdown() error
	SetTest(test test.Test)
}

type server struct {
	settings *Settings
	server   *http.Server
	test     test.Test
}

func NewWebServer(settings *Settings) Server {
	return &server{
		settings: settings,
	}
}

func (s *server) Listen() error {
	router := mux.NewRouter()

	router.HandleFunc("/gateways", s.gateways).Methods(http.MethodGet)
	router.PathPrefix("/").HandlerFunc(s.handler)

	s.server = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", host, s.settings.Port),
	}

	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func (s *server) SetTest(test test.Test) {
	s.test = test
}

func (s *server) gateways(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Result []*Gateway `json:"result"`
	}{
		Result: []*Gateway{
			{
				Url: fmt.Sprintf("http://%s:%d", host, s.settings.Port),
			},
		},
	}

	if bytes, err := json.Marshal(response); err != nil {
		test.InternalServerError(w)
	} else {
		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(bytes); err != nil {
			test.InternalServerError(w)
		}
	}
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	if s.test == nil {
		test.InternalServerError(w)
	} else if handler := s.test.Handler(); handler == nil {
		test.InternalServerError(w)
	} else {
		handler(w, r)
	}
}
