package batch_query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk    = "/batch_query"
	pathMock   = "/data_batch"
	query      = "tickets/resolve/allow"
	tenant     = "acmecorp"
	subject    = "alice"
	batchLimit = 20
	totalSize  = 90
)

type imap map[string]interface{}

func batchQuery() test.Test {
	apiRequestItems := make([]imap, totalSize)
	for i := 0; i < totalSize; i++ {
		apiRequestItems[i] = imap{
			"path": query,
		}
	}

	apiRequest := imap{
		"items": apiRequestItems,
		"input": imap{
			"tenant":  tenant,
			"subject": subject,
		},
	}

	apiResponseItems := make([]imap, totalSize)
	for i := 0; i < totalSize; i++ {
		apiResponseItems[i] = imap{
			"result": true,
		}
	}

	apiResponse := imap{
		"result": apiResponseItems,
	}

	cb := newCallbacks()

	settings := &test.Settings{
		Name: "batch-query",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodPost,
				Headers: test.DefaultContentTypeHeader,
				Cookies: test.AuthzCookie(tenant, subject),
				Body:    apiRequest,
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			pathMock: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPost),
					cb.Check,
				},
				Response: cb.Response,
			},
		},
	}

	return test.New(settings)
}

type input struct {
	Tenant  string `json:"tenant"`
	Subject string `json:"subject"`
}

type item struct {
	Path  string `json:"path"`
	Input *input `json:"input"`
}

type request struct {
	Items []*item `json:"items"`
	Input *input  `json:"input"`
}

type callbacks struct {
	Check    test.CheckRequest
	Response test.EmitResponse
}

// This sets up callbacks for checking batch requests and emitting
// responses. Among other things, it also checks the following:
//
// * That each batch is well formed.
// * That each batch respects the batch limit.
// * That the total size is not exceeded.
//
// Finally, it makes these checks in a thread-safe manner as multiple
// batch requests could be in flight at the same time.
func newCallbacks() *callbacks {
	var mutex sync.Mutex
	var total int

	return &callbacks{
		Check: func(w http.ResponseWriter, r *test.MockRequest) error {
			count, err := getCount(r)
			if err != nil {
				test.BadRequest(w)

				return err
			}

			if count > batchLimit {
				test.BadRequest(w)

				return fmt.Errorf("request: batch limit exceeded")
			}

			// Here, we must increment total in a thread-safe way.
			// We also must make a copy when performing the check.
			mutex.Lock()
			total += count
			myTotal := total
			mutex.Unlock()

			if myTotal > totalSize {
				test.BadRequest(w)

				return errors.New("request: exceeded total size")
			}

			input := &input{
				Tenant:  tenant,
				Subject: subject,
			}

			expected := &request{
				Items: make([]*item, count),
				Input: input,
			}

			for i := 0; i < count; i++ {
				expected.Items[i] = &item{
					Path:  query,
					Input: input,
				}
			}

			return test.CheckRequestBody(expected)(w, r)
		},
		Response: func(w http.ResponseWriter, r *test.MockRequest) error {
			count, err := getCount(r)
			if err != nil {
				return err
			}

			response := &struct {
				Result []imap `json:"result"`
			}{
				Result: make([]imap, count),
			}

			for i := 0; i < count; i++ {
				response.Result[i] = imap{
					"result": true,
				}
			}

			return test.DefaultResponse(200, response)(w, r)
		},
	}
}

func getCount(r *test.MockRequest) (int, error) {
	actual := &request{}

	if err := json.Unmarshal(r.Body, &actual); err != nil {
		return 0, err
	}

	return len(actual.Items), nil
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			batchQuery(),
		}
	}
}
