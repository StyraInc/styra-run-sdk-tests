package batch_query

import (
	"fmt"
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk    = "/"
	pathMock   = "/data_batch"
	query      = "tickets/resolve/allow"
	tenant     = "acmecorp"
	subject    = "alice"
	batchSize  = 30
	batchLimit = 20
)

type imap map[string]interface{}

func batchQuery() test.Test {
	apiRequestItems := make([]imap, 0)
	for i := 0; i < batchSize; i++ {
		apiRequestItems = append(
			apiRequestItems,
			imap{
				"path": query,
			},
		)
	}

	apiRequest := imap{
		"items": apiRequestItems,
		"input": imap{
			"tenant":  tenant,
			"subject": subject,
		},
	}

	apiResponseItems := make([]imap, 0)
	for i := 0; i < batchSize; i++ {
		apiResponseItems = append(
			apiResponseItems,
			imap{
				"result": true,
			},
		)
	}

	apiResponse := imap{
		"result": apiResponseItems,
	}

	cb := newCallbacks(batchSize)

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

// This works by emitting the check and response functions
// as closures that close over remaining and count. This
// allows the callbacks to track the remaining number of
// batches to check for and emit.
func newCallbacks(batchSize int) *callbacks {
	remaining := batchSize
	count := 0

	return &callbacks{
		Check: func(w http.ResponseWriter, r *http.Request) error {
			if remaining <= 0 {
				test.BadRequest(w)

				return fmt.Errorf("request: batch count exceeded")
			}

			count = remaining
			if count > batchLimit {
				count = batchLimit
			}
			remaining -= count

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
		Response: func(w http.ResponseWriter, r *http.Request) error {
			responseItems := make([]imap, 0)
			for i := 0; i < count; i++ {
				responseItems = append(
					responseItems,
					imap{
						"result": true,
					},
				)
			}

			response := &struct {
				Result []imap `json:"result"`
			}{
				Result: responseItems,
			}

			return test.DefaultResponse(200, response)(w, r)
		},
	}
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			batchQuery(),
		}
	}
}
