package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/util"
)

type Request struct {
	Path    string
	Method  string
	Headers map[string]string
	Cookies []*http.Cookie
	Body    interface{}
}

type Response struct {
	Rest *util.Rest
	Body interface{}
}

type CheckResponse func(response *Response) error

type Api struct {
	Request *Request
	Checks  []CheckResponse
}

type CheckRequest func(w http.ResponseWriter, r *http.Request) error
type EmitResponse func(w http.ResponseWriter, r *http.Request) error

type Mock struct {
	Checks   []CheckRequest
	Response EmitResponse
}

type Settings struct {
	Name  string
	Api   *Api
	Mocks map[string]*Mock
}

type test struct {
	settings *Settings
	errors   []error
}

func New(settings *Settings) Test {
	return &test{
		settings: settings,
		errors:   make([]error, 0),
	}
}

func (t *test) Name() string {
	return t.settings.Name
}

func (t *test) Run(host string) []error {
	request := t.settings.Api.Request

	var body interface{}
	rest := &util.Rest{
		Url:     fmt.Sprintf("%s%s", host, request.Path),
		Method:  request.Method,
		Cookies: request.Cookies,
		Headers: request.Headers,
		Decoder: util.JsonDecoder(&body),
	}

	if request.Body != nil {
		rest.Encoder = util.JsonEncoder(request.Body)
	}

	if err := rest.Execute(context.Background()); err != nil {
		t.errors = append(t.errors, err)
	} else {
		response := &Response{
			Rest: rest,
			Body: body,
		}

		for _, check := range t.settings.Api.Checks {
			if err := check(response); err != nil {
				t.errors = append(t.errors, err)
			}
		}
	}

	return t.errors
}

func (t *test) Handler() http.HandlerFunc {
	return t.handler
}

func (t *test) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()

	if mock, ok := t.settings.Mocks[path]; ok {
		for _, check := range mock.Checks {
			if err := check(w, r); err != nil {
				t.errors = append(t.errors, err)

				return
			}
		}

		if mock.Response == nil {
			t.errors = append(t.errors, t.missingBodyCallbackError())

			InternalServerError(w)
		} else if err := mock.Response(w, r); err != nil {
			t.errors = append(t.errors, err)
		}
	} else {
		t.errors = append(t.errors, t.unexpectedRequestError(path))

		InternalServerError(w)
	}
}

func (t *test) missingBodyCallbackError() error {
	return fmt.Errorf("missing body callback")
}

func (t *test) unexpectedRequestError(path string) error {
	return fmt.Errorf("unexpected request: %s", path)
}
