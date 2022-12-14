package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/styrainc/styra-run-sdk-tests/util"
)

type Request struct {
	Path    string
	Method  string
	Cookies []*http.Cookie
	Headers map[string]string
	Queries map[string]string
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

type MockRequest struct {
	Request *http.Request
	Body    []byte
}

type CheckRequest func(w http.ResponseWriter, r *MockRequest) error
type EmitResponse func(w http.ResponseWriter, r *MockRequest) error

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
	mutex    sync.Mutex
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
		Queries: request.Queries,
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

// Here, errors must be accumulated in a thread-safe manner as
// multiple requests may be received in parallel in some cases.
func (t *test) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		t.addError(err)

		InternalServerError(w)

		return
	}

	mockRequest := &MockRequest{
		Request: r,
		Body:    bytes,
	}

	if mock, ok := t.settings.Mocks[path]; ok {
		for _, check := range mock.Checks {
			if err := check(w, mockRequest); err != nil {
				t.addError(err)

				return
			}
		}

		if mock.Response == nil {
			t.addError(t.missingResponseCallbackError())

			InternalServerError(w)
		} else if err := mock.Response(w, mockRequest); err != nil {
			t.addError(err)
		}
	} else {
		t.addError(t.unexpectedRequestError(path))

		InternalServerError(w)
	}
}

func (t *test) addError(err error) {
	t.mutex.Lock()
	t.errors = append(t.errors, err)
	t.mutex.Unlock()
}

func (t *test) missingResponseCallbackError() error {
	return fmt.Errorf("missing response callback")
}

func (t *test) unexpectedRequestError(path string) error {
	return fmt.Errorf("unexpected request: %s", path)
}
