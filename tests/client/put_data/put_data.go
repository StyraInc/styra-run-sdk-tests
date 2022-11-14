package put_data

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

func putData() test.Test {
	apiRequest := imap{
		"alice": ilist{
			"ADMIN",
		},
		"billy": ilist{
			"VIEWER",
		},
		"bob": ilist{
			"VIEWER",
		},
	}

	apiResponse := imap{}

	mockRequest := apiRequest

	mockResponse := imap{}

	settings := &test.Settings{
		Name: "put-data",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodPut,
				Headers: test.DefaultContentTypeHeader,
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
					test.CheckRequestMethod(http.MethodPut),
					test.CheckRequestContentType(test.ApplicationJson),
					test.CheckRequestBody(mockRequest),
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
			putData(),
		}
	}
}
