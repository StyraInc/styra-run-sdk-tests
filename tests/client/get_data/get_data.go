package get_data

import (
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/data/rbac/user_bindings/acmecorp"
	pathMock = "/data/rbac/user_bindings/acmecorp"
)

type imap map[string]interface{}
type ilist []interface{}

func getData() test.Test {
	apiResponse := imap{
		"result": imap{
			"alice": ilist{
				"ADMIN",
			},
			"billy": ilist{
				"VIEWER",
			},
			"bob": ilist{
				"VIEWER",
			},
		},
	}

	mockResponse := apiResponse

	settings := &test.Settings{
		Name: "get-data",
		Api: &test.Api{
			Request: &test.Request{
				Path:   pathSdk,
				Method: http.MethodGet,
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			pathMock: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodGet),
				},
				Response: test.DefaultResponse(200, mockResponse),
			},
		},
	}

	return test.New(settings)
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			getData(),
		}
	}
}
