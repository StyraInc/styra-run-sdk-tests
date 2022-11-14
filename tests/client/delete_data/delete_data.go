package delete_data

import (
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/data/rbac/user_bindings/acmecorp"
	pathMock = "/data/rbac/user_bindings/acmecorp"
)

type imap map[string]interface{}

func deleteData() test.Test {
	apiResponse := imap{}

	mockResponse := imap{}

	settings := &test.Settings{
		Name: "delete-data",
		Api: &test.Api{
			Request: &test.Request{
				Path:   pathSdk,
				Method: http.MethodDelete,
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			pathMock: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodDelete),
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
			deleteData(),
		}
	}
}
