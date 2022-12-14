package get_roles

import (
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/roles"
	pathMock = "/data/rbac/roles"
	tenant   = "acmecorp"
	subject  = "alice"
)

type imap map[string]interface{}
type ilist []interface{}

func getRoles() test.Test {
	apiResponse := imap{
		"result": ilist{
			"ADMIN",
			"VIEWER",
		},
	}

	mockResponse := imap{
		"result": ilist{
			"ADMIN",
			"VIEWER",
		},
	}

	settings := &test.Settings{
		Name: "get-roles",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodGet,
				Cookies: test.AuthzCookie(tenant, subject),
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			test.AuthzPath: test.AuthzMock(tenant, subject, true),
			pathMock: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPost),
					test.CheckRequestContentType(test.ApplicationJson),
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
			getRoles(),
		}
	}
}
